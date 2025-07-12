# LogHarbor Processor

The core service of the LogHarbor distributed logging system. Consumes log messages from Apache Kafka, indexes them to Elasticsearch, and provides real-time streaming via WebSocket.

## Features

- 📥 **Kafka Consumer**: Reliable message consumption with proper error handling
- 🔍 **Elasticsearch Indexing**: Structured log storage with metadata
- �� **WebSocket Server**: Real-time log streaming to connected clients
- �� **Health Checks**: Built-in health monitoring endpoint
- 🔧 **Configurable**: Environment-based configuration
- 🐳 **Docker Ready**: Containerized with multi-stage build
- 🛡️ **Production Ready**: Graceful shutdown, connection management, and error handling

## Architecture




## 4. Create a configuration file

```yaml:log-processor/config.yaml
# LogHarbor Processor Configuration
kafka:
  bootstrap_servers: "localhost:9092"
  topic: "logs"
  group_id: "log-processor-group"
  auto_offset_reset: "earliest"
  enable_auto_commit: true
  auto_commit_interval_ms: 1000
  session_timeout_ms: 30000
  heartbeat_interval_ms: 3000

elasticsearch:
  url: "http://localhost:9200"
  index: "logs"
  sniff: false
  healthcheck: false

websocket:
  port: "8080"
  read_buffer_size: 1024
  write_buffer_size: 1024

health:
  port: "8081"
```

## 5. Create a README for the log-processor

```markdown:log-processor/README.md
# LogHarbor Processor

The core service of the LogHarbor distributed logging system. Consumes log messages from Apache Kafka, indexes them to Elasticsearch, and provides real-time streaming via WebSocket.

## Features

- 📥 **Kafka Consumer**: Reliable message consumption with proper error handling
- 🔍 **Elasticsearch Indexing**: Structured log storage with metadata
-  **WebSocket Server**: Real-time log streaming to connected clients
-  **Health Checks**: Built-in health monitoring endpoint
- 🔧 **Configurable**: Environment-based configuration
- 🐳 **Docker Ready**: Containerized with multi-stage build
- 🛡️ **Production Ready**: Graceful shutdown, connection management, and error handling

## Architecture

```
Kafka → Log Processor → Elasticsearch (Storage)
                ↓
            WebSocket → Dashboard (Real-time)
```
<code_block_to_apply_changes_from>
```
log-processor/
├── main.go          # Main application entry point
├── go.mod           # Go module definition
├── Dockerfile       # Container definition
├── config.yaml      # Configuration file
└── README.md        # This file
```

### Building
```bash
go build -o log-processor .
```

### Testing
```bash
go test ./...
```

## Monitoring

### Health Check Response
```json
{
  "status": "healthy",
  "timestamp": "2024-01-15T10:30:00Z",
  "service": "log-processor",
  "version": "1.0.0",
  "clients": 5
}
```

### Log Levels
- `` - Service startup
- `📥` - Kafka consumer events
- `📨` - Message consumption
- `✅` - Successful operations
- `❌` - Errors
- `⚠️` - Warnings
- `🔗` - WebSocket connections
- `🔌` - WebSocket disconnections
- `🏥` - Health check events
- `` - Shutdown events

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

MIT License - see LICENSE file for details
```

## How to Run the Enhanced Log Processor

1. **Initialize the Go module:**
   ```bash
   cd log-processor
   go mod init log-processor
   go mod tidy
   ```

2. **Run locally:**
   ```bash
   go run main.go
   ```

3. **Run with Docker:**
   ```bash
   docker build -t log-harbor-processor .
   docker run -p 8080:8080 -p 8081:8081 \
     -e KAFKA_BOOTSTRAP_SERVERS=kafka:29092 \
     -e ELASTICSEARCH_URL=http://elasticsearch:9200 \
     log-harbor-processor
   ```

## Key Improvements Made

✅ **Better Architecture**: Clean separation of concerns with dedicated types and methods

✅ **Robust Error Handling**: Comprehensive error handling for Kafka, Elasticsearch, and WebSocket operations

✅ **Graceful Shutdown**: Proper cleanup and shutdown handling

✅ **Connection Management**: Better WebSocket client management with read/write pumps

✅ **Health Monitoring**: Built-in health check endpoint

✅ **Configuration Management**: Environment-based configuration

✅ **Production Features**: Connection pooling, timeouts, and retry logic

✅ **Structured Logging**: Better logging with emojis and structured data

✅ **Docker Support**: Multi-stage Docker build for production deployment

## Next Steps

Now that we have both backend services (log-agent and log-processor) enhanced, the next logical step would be to:

1. **Test the backend services together** - Make sure they communicate properly
2. **Enhance the React frontend** - Improve the dashboard with better UI/UX
3. **Update the Docker Compose file** - Add all services to the compose file
4. **Add search functionality** - Implement Elasticsearch search queries in the frontend

Would you like to test the backend services first, or should we move on to enhancing the React frontend? 