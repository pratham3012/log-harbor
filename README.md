# LogHarbor 🚢

LogHarbor is a distributed log ingestion pipeline designed to reliably collect, process, and visualize high-volume log data in real-time. This project demonstrates a modern microservices architecture using Go, Kafka, Elasticsearch, and React.

## Features

-   **Distributed Log Collection**: Go-based agents collect logs from multiple sources.
-   **Reliable Message Queuing**: Kafka provides a durable and scalable buffer for incoming logs.
-   **Efficient Log Processing**: A Go backend service consumes logs from Kafka, parses them, and indexes them into Elasticsearch.
-   **Real-time Visualization**: A React dashboard with WebSockets displays a live stream of logs.
-   **Containerized**: The entire stack is containerized with Docker for easy setup and portability.

## Architecture

![Architecture Diagram](https://i.imgur.com/7g5u2s1.png)

## Tech Stack

-   **Backend**: Go
-   **Message Queue**: Apache Kafka
-   **Search & Indexing**: Elasticsearch
-   **Frontend**: React
-   **Containerization**: Docker, Docker Compose

## Getting Started

### Prerequisites

-   Go (v1.18+)
-   Node.js (v16+)
-   Docker & Docker Compose

### 1. Clone the Repository

```bash
git clone [https://github.com/your-username/logharbor.git](https://github.com/your-username/logharbor.git)
cd logharbor
```

### 2. Start Infrastructure Services

This will start Kafka and Elasticsearch in Docker containers.

```bash
docker-compose up -d
```

### 3. Run the Backend Services

-   **Log Processor (Consumer & WebSocket Server)**

    ```bash
    cd log-processor
    go run main.go
    ```

-   **Log Agent (Producer)**

    ```bash
    cd log-agent
    go run main.go
    ```

### 4. Run the Frontend

```bash
cd log-dashboard
npm install
npm start
```

 Access Your System:
📊 Dashboard: http://localhost:3000
🏥 Agent Health: http://localhost:8080/health
🏥 Processor Health: http://localhost:8081/health
🔍 Elasticsearch: http://localhost:9200

Logs are being generated every 2 seconds by the agent
Processed and indexed to Elasticsearch via the processor
Real-time WebSocket streaming available for the dashboard
Complete observability pipeline is working!


Open your browser to `http://localhost:3000` to see the live log dashboard.

## Deployment

This project is designed for easy deployment. See the official documentation for services like [Render](https://render.com/), [Fly.io](https://fly.io/), or cloud providers like AWS/GCP for deploying Docker containers and managed Kafka/Elasticsearch services.