package main

import (
    "encoding/json"
    "fmt"
    "log"
    "math/rand"
    "os"
    "os/signal"
    "syscall"
    "time"

    "github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

// LogEntry represents a structured log message
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

// LogGenerator handles log message generation
type LogGenerator struct {
    producer *kafka.Producer
    topic    string
    service  string
}

// NewLogGenerator creates a new log generator instance
func NewLogGenerator(bootstrapServers, topic, service string) (*LogGenerator, error) {
    producer, err := kafka.NewProducer(&kafka.ConfigMap{
        "bootstrap.servers": bootstrapServers,
        "client.id":         fmt.Sprintf("log-agent-%s", service),
        "acks":              "all",
        "retries":           3,
    })
    if err != nil {
        return nil, fmt.Errorf("failed to create producer: %s", err)
    }

    return &LogGenerator{
        producer: producer,
        topic:    topic,
        service:  service,
    }, nil
}

// Close closes the producer
func (lg *LogGenerator) Close() {
    lg.producer.Close()
}

// generateLogEntry creates a realistic log entry
func (lg *LogGenerator) generateLogEntry() LogEntry {
    logLevels := []string{"INFO", "WARN", "ERROR", "DEBUG"}
    level := logLevels[rand.Intn(len(logLevels))]

    messages := lg.getMessagesForLevel(level)
    message := messages[rand.Intn(len(messages))]

    logEntry := LogEntry{
        Level:     level,
        Message:   message,
        Timestamp: time.Now().Format(time.RFC3339),
        Service:   lg.service,
    }

    // Add additional fields for certain log levels
    switch level {
    case "INFO", "ERROR":
        logEntry.UserID = lg.getRandomUserID()
        logEntry.RequestID = lg.getRandomRequestID()
        if level == "ERROR" {
            logEntry.Duration = rand.Int63n(5000) + 100 // 100-5100ms
        }
    case "DEBUG":
        logEntry.Duration = rand.Int63n(200) + 10 // 10-210ms
    }

    // Add IP for some logs
    if rand.Float32() < 0.7 { // 70% chance
        logEntry.IP = lg.getRandomIP()
    }

    return logEntry
}

// getMessagesForLevel returns appropriate messages for each log level
func (lg *LogGenerator) getMessagesForLevel(level string) []string {
    messages := map[string][]string{
        "INFO": {
            "User logged in successfully",
            "Payment processed successfully",
            "User profile updated",
            "Email notification sent",
            "Database connection established",
            "Cache hit for user data",
            "File uploaded successfully",
            "API request completed",
            "User session created",
            "Data backup completed",
        },
        "WARN": {
            "High memory usage detected",
            "Database connection pool at 80% capacity",
            "Slow query detected (>2s)",
            "Rate limit approaching threshold",
            "Deprecated API endpoint called",
            "Large file upload detected",
            "Unusual login pattern detected",
            "Cache miss rate increased",
            "External service response time degraded",
            "Disk space usage at 85%",
        },
        "ERROR": {
            "Failed to connect to database",
            "Payment gateway timeout",
            "Invalid authentication token",
            "File upload failed - disk full",
            "External API call failed",
            "Database transaction rollback",
            "User session expired",
            "Email delivery failed",
            "Cache connection lost",
            "SSL certificate expired",
        },
        "DEBUG": {
            "Processing user request",
            "Validating input parameters",
            "Executing database query",
            "Sending HTTP request",
            "Parsing response data",
            "Checking user permissions",
            "Generating authentication token",
            "Compressing response data",
            "Logging audit trail",
            "Updating cache entry",
        },
    }
    return messages[level]
}

// getRandomUserID returns a random user ID
func (lg *LogGenerator) getRandomUserID() string {
    userIDs := []string{"user123", "user456", "user789", "user101", "user202", "admin001", "guest999"}
    return userIDs[rand.Intn(len(userIDs))]
}

// getRandomRequestID returns a random request ID
func (lg *LogGenerator) getRandomRequestID() string {
    requestIDs := []string{"req-001", "req-002", "req-003", "req-004", "req-005", "req-006", "req-007"}
    return requestIDs[rand.Intn(len(requestIDs))]
}

// getRandomIP returns a random IP address
func (lg *LogGenerator) getRandomIP() string {
    ips := []string{"192.168.1.100", "10.0.0.50", "172.16.0.25", "203.0.113.10", "198.51.100.5"}
    return ips[rand.Intn(len(ips))]
}

// produceLog sends a log entry to Kafka
func (lg *LogGenerator) produceLog(logEntry LogEntry) error {
    jsonData, err := json.Marshal(logEntry)
    if err != nil {
        return fmt.Errorf("error marshaling log entry: %s", err)
    }

    err = lg.producer.Produce(&kafka.Message{
        TopicPartition: kafka.TopicPartition{Topic: &lg.topic, Partition: kafka.PartitionAny},
        Value:          jsonData,
        Key:            []byte(logEntry.Service),
    }, nil)

    if err != nil {
        return fmt.Errorf("error producing message: %s", err)
    }

    // Flush to ensure message is sent
    lg.producer.Flush(1000)
    return nil
}

// Start begins the log generation process
func (lg *LogGenerator) Start() {
    log.Printf(" Starting LogHarbor Agent for service: %s", lg.service)
    log.Printf("📡 Producing logs to topic: %s", lg.topic)

    // Set up signal handling for graceful shutdown
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

    ticker := time.NewTicker(2 * time.Second)
    defer ticker.Stop()

    for {
        select {
        case <-ticker.C:
            logEntry := lg.generateLogEntry()
            
            if err := lg.produceLog(logEntry); err != nil {
                log.Printf("❌ Error producing log: %s", err)
                continue
            }

            jsonData, _ := json.Marshal(logEntry)
            fmt.Printf("✅ Produced log: %s\n", string(jsonData))

        case <-sigChan:
            log.Println("🛑 Received shutdown signal, stopping log agent...")
            return
        }
    }
}

func main() {
    // Configuration
    bootstrapServers := os.Getenv("KAFKA_BOOTSTRAP_SERVERS")
    if bootstrapServers == "" {
        bootstrapServers = "localhost:9092"
    }

    topic := os.Getenv("KAFKA_TOPIC")
    if topic == "" {
        topic = "logs"
    }

    service := os.Getenv("SERVICE_NAME")
    if service == "" {
        service = "log-agent"
    }

    startHealthServer("8080", service)

    // Create log generator
    generator, err := NewLogGenerator(bootstrapServers, topic, service)
    if err != nil {
        log.Fatalf("❌ Failed to create log generator: %s", err)
    }
    defer generator.Close()

    // Start generating logs
    generator.Start()
}