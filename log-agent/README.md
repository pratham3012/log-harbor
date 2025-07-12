# LogHarbor Agent

A lightweight Go service that generates realistic log messages and sends them to Apache Kafka for the LogHarbor distributed logging system.

## Features

- 🚀 **Realistic Log Generation**: Produces structured log messages with different levels (INFO, WARN, ERROR, DEBUG)
- 📡 **Kafka Integration**: Sends logs to Apache Kafka for reliable message queuing
- 🔧 **Configurable**: Environment-based configuration for different deployment scenarios
- 🐳 **Docker Ready**: Containerized with multi-stage Docker build
- �� **Health Checks**: Built-in health check endpoint
- 🛡️ **Production Ready**: Proper error handling, graceful shutdown, and logging

## Quick Start

### Local Development

1. **Prerequisites**
   - Go 1.21+
   - Apache Kafka running on localhost:9092

2. **Run the agent**
   ```bash
   go mod tidy
   go run main.go
   ```

### With Docker

```bash
docker build -t log-harbor-agent .
docker run -e KAFKA_BOOTSTRAP_SERVERS=kafka:29092 log-harbor-agent
```

### Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `KAFKA_BOOTSTRAP_SERVERS` | `localhost:9092` | Kafka broker addresses |
| `KAFKA_TOPIC` | `logs` | Kafka topic name |
| `SERVICE_NAME` | `log-agent` | Service identifier |

## Log Message Format

```json
{
  "level": "INFO",
  "message": "User logged in successfully",
  "timestamp": "2024-01-15T10:30:00Z",
  "service": "auth-service",
  "user_id": "user123",
  "request_id": "req-001",
  "ip": "192.168.1.100",
  "duration_ms": 150
}
```

## Architecture

The LogHarbor Agent is part of a microservices architecture:

```
Log Agent → Kafka → Log Processor → Elasticsearch + WebSocket → Dashboard
```

## Development

### Project Structure
```
log-agent/
├── main.go          # Main application entry point
├── health.go        # Health check endpoint
├── go.mod           # Go module definition
├── Dockerfile       # Container definition
├── config.yaml      # Configuration file
└── README.md        # This file
```

## How to Run the Enhanced Log Agent

1. **Initialize the Go module:**
   ```bash
   cd log-agent
   go mod init log-agent
   go mod tidy
   ```

2. **Run locally:**
   ```bash
   go run main.go
   ```

3. **Run with Docker:**
   ```bash
   docker build -t log-harbor-agent .
   docker run -e KAFKA_BOOTSTRAP_SERVERS=kafka:29092 log-harbor-agent
   ```

## Key Improvements Made

✅ **Structured Logging**: JSON format with additional fields like `service`, `user_id`, `request_id`, `ip`, and `duration_ms`

✅ **Better Error Handling**: Proper error handling and graceful shutdown

✅ **Environment Configuration**: Configurable via environment variables

✅ **Production Ready**: Docker support, health checks, and proper logging

✅ **Realistic Data**: More diverse and realistic log messages for different scenarios

✅ **Clean Architecture**: Separation of concerns with dedicated types and methods

✅ **Documentation**: Comprehensive README and inline comments

The log-agent is now a robust, production-ready service that generates realistic log data and sends it to Kafka. It's perfect for demonstrating your backend development skills with proper error handling, configuration management, and containerization! 

Would you like me to continue with the log-processor service next, or would you like to test this enhanced log-agent first?

### Building
```bash
go build -o log-agent .
```

### Testing
```bash
go test ./...
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

MIT License - see LICENSE file for details