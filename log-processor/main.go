package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/gorilla/websocket"
	"github.com/olivere/elastic/v7"
)

// LogEntry represents the structure of log messages
type LogEntry struct {
	Level     string `json:"level"`
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
	Service   string `json:"service"`
	UserID    string `json:"user_id,omitempty"`
	RequestID string `json:"request_id,omitempty"`
	IP        string `json:"ip,omitempty"`
	Duration  int64  `json:"duration_ms,omitempty"`
}

// WebSocketClient represents a connected WebSocket client
type WebSocketClient struct {
	conn      *websocket.Conn
	send      chan []byte
	processor *LogProcessor
}

// LogProcessor handles log processing, indexing, and WebSocket broadcasting
type LogProcessor struct {
	esClient   *elastic.Client
	consumer   *kafka.Consumer
	clients    map[*WebSocketClient]bool
	broadcast  chan []byte
	mutex      sync.RWMutex
	ctx        context.Context
	cancel     context.CancelFunc
	config     *Config
}

// Config holds application configuration
type Config struct {
	KafkaBootstrapServers string
	KafkaTopic           string
	KafkaGroupID         string
	ElasticsearchURL     string
	ElasticsearchIndex   string
	WebSocketPort        string
	HealthPort           string
}

// NewConfig creates a new configuration from environment variables
func NewConfig() *Config {
	return &Config{
		KafkaBootstrapServers: getEnv("KAFKA_BOOTSTRAP_SERVERS", "localhost:9092"),
		KafkaTopic:           getEnv("KAFKA_TOPIC", "logs"),
		KafkaGroupID:         getEnv("KAFKA_GROUP_ID", "log-processor-group"),
		ElasticsearchURL:     getEnv("ELASTICSEARCH_URL", "http://localhost:9200"),
		ElasticsearchIndex:   getEnv("ELASTICSEARCH_INDEX", "logs"),
		WebSocketPort:        getEnv("WS_PORT", "8080"),
		HealthPort:           getEnv("HEALTH_PORT", "8081"),
	}
}

// getEnv gets an environment variable with a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// NewLogProcessor creates a new log processor instance
func NewLogProcessor(config *Config) (*LogProcessor, error) {
	ctx, cancel := context.WithCancel(context.Background())

	// Initialize Elasticsearch client with correct v7 function names
	esClient, err := elastic.NewClient(
		elastic.SetURL(config.ElasticsearchURL),
		elastic.SetSniff(false),
		elastic.SetHealthcheck(false),
		elastic.SetRetrier(elastic.NewBackoffRetrier(elastic.NewExponentialBackoff(100*time.Millisecond, 30*time.Second))),
	)
	if err != nil {
		cancel()
		return nil, fmt.Errorf("failed to create Elasticsearch client: %w", err)
	}

	// Initialize Kafka consumer
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": config.KafkaBootstrapServers,
		"group.id":          config.KafkaGroupID,
		"auto.offset.reset": "earliest",
		"enable.auto.commit": true,
		"auto.commit.interval.ms": 1000,
		"session.timeout.ms":      30000,
		"heartbeat.interval.ms":   3000,
	})
	if err != nil {
		cancel()
		return nil, fmt.Errorf("failed to create Kafka consumer: %w", err)
	}

	// Subscribe to topic
	err = consumer.SubscribeTopics([]string{config.KafkaTopic}, nil)
	if err != nil {
		cancel()
		return nil, fmt.Errorf("failed to subscribe to topic: %w", err)
	}

	return &LogProcessor{
		esClient:   esClient,
		consumer:   consumer,
		clients:    make(map[*WebSocketClient]bool),
		broadcast:  make(chan []byte, 1000),
		ctx:        ctx,
		cancel:     cancel,
		config:     config,
	}, nil
}

// Start begins the log processing
func (lp *LogProcessor) Start() error {
	log.Printf(" Starting LogHarbor Processor")
	log.Printf("📡 Kafka: %s (topic: %s)", lp.config.KafkaBootstrapServers, lp.config.KafkaTopic)
	log.Printf("🔍 Elasticsearch: %s (index: %s)", lp.config.ElasticsearchURL, lp.config.ElasticsearchIndex)
	log.Printf("🌐 WebSocket: :%s", lp.config.WebSocketPort)

	// Start goroutines
	go lp.consumeAndIndex()
	go lp.handleBroadcasts()
	go lp.startHealthServer()
	go lp.startWebSocketServer()

	// Wait for shutdown signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("🛑 Received shutdown signal, stopping processor...")
	return lp.Shutdown()
}

// Shutdown gracefully shuts down the processor
func (lp *LogProcessor) Shutdown() error {
	log.Println(" Shutting down LogHarbor Processor...")

	// Cancel context
	lp.cancel()

	// Close Kafka consumer
	if lp.consumer != nil {
		lp.consumer.Close()
	}

	// Close all WebSocket connections
	lp.mutex.Lock()
	for client := range lp.clients {
		client.conn.Close()
		close(client.send)
	}
	lp.clients = make(map[*WebSocketClient]bool)
	lp.mutex.Unlock()

	// Close broadcast channel
	close(lp.broadcast)

	log.Println("✅ LogHarbor Processor shutdown complete")
	return nil
}

// consumeAndIndex consumes messages from Kafka and indexes them to Elasticsearch
func (lp *LogProcessor) consumeAndIndex() {
	log.Println("📥 Starting Kafka consumer...")

	for {
		select {
		case <-lp.ctx.Done():
			log.Println("📥 Kafka consumer stopped")
			return
		default:
			msg, err := lp.consumer.ReadMessage(100 * time.Millisecond)
			if err != nil {
				if kafkaErr, ok := err.(kafka.Error); ok && kafkaErr.Code() == kafka.ErrTimedOut {
					continue
				}
				log.Printf("❌ Kafka consumer error: %v", err)
				continue
			}

			log.Printf("📨 Consumed message from topic %s: %s", *msg.TopicPartition.Topic, string(msg.Value))

			// Index to Elasticsearch
			go lp.indexToElasticsearch(msg.Value)

			// Broadcast to WebSocket clients
			select {
			case lp.broadcast <- msg.Value:
			default:
				log.Println("⚠️ Broadcast channel full, dropping message")
			}
		}
	}
}

// indexToElasticsearch indexes a log message to Elasticsearch
func (lp *LogProcessor) indexToElasticsearch(data []byte) {
	// Parse the log entry to add metadata
	var logEntry LogEntry
	if err := json.Unmarshal(data, &logEntry); err != nil {
		log.Printf("❌ Error parsing log entry: %v", err)
		return
	}

	// Add indexing timestamp
	indexData := map[string]interface{}{
		"@timestamp": time.Now().Format(time.RFC3339),
		"log_entry":  logEntry,
	}

	// Index to Elasticsearch
	_, err := lp.esClient.Index().
		Index(lp.config.ElasticsearchIndex).
		BodyJson(indexData).
		Do(lp.ctx)

	if err != nil {
		log.Printf("❌ Error indexing to Elasticsearch: %v", err)
		return
	}

	log.Printf("✅ Indexed log to Elasticsearch: %s", logEntry.Message)
}

// handleBroadcasts handles broadcasting messages to WebSocket clients
func (lp *LogProcessor) handleBroadcasts() {
	for {
		select {
		case msg, ok := <-lp.broadcast:
			if !ok {
				return
			}
			lp.broadcastToClients(msg)
		case <-lp.ctx.Done():
			return
		}
	}
}

// broadcastToClients sends a message to all connected WebSocket clients
func (lp *LogProcessor) broadcastToClients(message []byte) {
	lp.mutex.RLock()
	defer lp.mutex.RUnlock()

	for client := range lp.clients {
		select {
		case client.send <- message:
		default:
			log.Printf("⚠️ Client send buffer full, removing client")
			client.conn.Close()
			delete(lp.clients, client)
		}
	}
}

// startWebSocketServer starts the WebSocket server
func (lp *LogProcessor) startWebSocketServer() {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true // Allow all connections for development
		},
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		lp.handleWebSocketConnection(w, r, upgrader)
	})

	log.Printf("🌐 WebSocket server starting on :%s", lp.config.WebSocketPort)
	if err := http.ListenAndServe(":"+lp.config.WebSocketPort, nil); err != nil {
		log.Printf("❌ WebSocket server error: %v", err)
	}
}

// handleWebSocketConnection handles a new WebSocket connection
func (lp *LogProcessor) handleWebSocketConnection(w http.ResponseWriter, r *http.Request, upgrader websocket.Upgrader) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("❌ WebSocket upgrade error: %v", err)
		return
	}

	client := &WebSocketClient{
		conn:      conn,
		send:      make(chan []byte, 256),
		processor: lp,
	}

	// Register client
	lp.mutex.Lock()
	lp.clients[client] = true
	clientCount := len(lp.clients)
	lp.mutex.Unlock()

	log.Printf("🔗 New WebSocket client connected. Total clients: %d", clientCount)

	// Send welcome message
	welcomeMsg := LogEntry{
		Level:     "INFO",
		Message:   "Connected to LogHarbor WebSocket",
		Timestamp: time.Now().Format(time.RFC3339),
		Service:   "log-processor",
	}
	welcomeJSON, _ := json.Marshal(welcomeMsg)
	client.send <- welcomeJSON

	// Start client handlers
	go client.writePump()
	go client.readPump()
}

// writePump pumps messages from the send channel to the WebSocket connection
func (client *WebSocketClient) writePump() {
	ticker := time.NewTicker(54 * time.Second)
	defer func() {
		ticker.Stop()
		client.conn.Close()
	}()

	for {
		select {
		case message, ok := <-client.send:
			if !ok {
				client.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := client.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			if err := client.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// readPump pumps messages from the WebSocket connection to the processor
func (client *WebSocketClient) readPump() {
	defer func() {
		client.processor.mutex.Lock()
		delete(client.processor.clients, client)
		clientCount := len(client.processor.clients)
		client.processor.mutex.Unlock()
		log.Printf("🔌 WebSocket client disconnected. Total clients: %d", clientCount)
		client.conn.Close()
	}()

	client.conn.SetReadLimit(512)
	client.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	client.conn.SetPongHandler(func(string) error {
		client.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		_, _, err := client.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("❌ WebSocket read error: %v", err)
			}
			break
		}
	}
}

// startHealthServer starts the health check server
func (lp *LogProcessor) startHealthServer() {
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		status := map[string]interface{}{
			"status":    "healthy",
			"timestamp": time.Now().Format(time.RFC3339),
			"service":   "log-processor",
			"version":   "1.0.0",
			"clients":   len(lp.clients),
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(status)
	})

	go func() {
		log.Printf("🏥 Health check server starting on :%s", lp.config.HealthPort)
		if err := http.ListenAndServe(":"+lp.config.HealthPort, nil); err != nil {
			log.Printf("❌ Health server error: %v", err)
		}
	}()
}

func main() {
	// Initialize configuration
	config := NewConfig()

	// Create and start log processor
	processor, err := NewLogProcessor(config)
	if err != nil {
		log.Fatalf("❌ Failed to create log processor: %v", err)
	}

	// Start the processor
	if err := processor.Start(); err != nil {
		log.Fatalf("❌ Failed to start log processor: %v", err)
	}
}