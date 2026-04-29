# StreamPulse Ingestion API

StreamPulse is a high-performance, distributed event ingestion system designed to handle high-load telemetry and event data. The project demonstrates a modern microservices architecture using Go, RabbitMQ, PostgreSQL, and Redis.

## 🚀 Project Goal
The primary objective of this project is to build a scalable pipeline that can accept a massive stream of events (clicks, views, transactions) via a REST API, buffer them using a message broker to prevent system overload, and process them asynchronously before storing them in a persistent database.

## 🏗 Architecture
The system is split into two main components:
1. **Ingestion API**: A lightweight Go service that validates incoming JSON events and immediately pushes them into RabbitMQ. It is designed for maximum throughput and low latency.
2. **Processor Worker**: A background consumer that retrieves events from RabbitMQ, enriches them with real-world data (like GeoIP via external APIs), ensures data integrity (idempotency check via Redis), and saves them to PostgreSQL.

### Tech Stack
- **Language**: Go (Golang) 1.25+
- **Message Broker**: RabbitMQ (Asynchronous buffering)
- **Database**: PostgreSQL (Relational storage)
- **Cache**: Redis (Idempotency / Deduplication)
- **Containerization**: Docker & Docker Compose

## 🛠 Setup & Installation

### Prerequisites
- Docker and Docker Compose installed on your machine.
- `curl` or Postman for testing.

### Running with Docker
The entire infrastructure and the application services can be started with a single command:

```bash
docker compose up -d --build
```

This will spin up:
- **sp_api**: The Ingestion API at http://localhost:8080
- **sp_worker**: The background event processor.
- **sp_postgres**: PostgreSQL database on port 5432.
- **sp_rabbitmq**: RabbitMQ management UI at http://localhost:15672 (guest/guest).
- **sp_redis**: Redis instance on port 6379.

## 📡 API Usage
### Send an Event
To send a telemetry event, use the following endpoint:
`POST /v1/events`

Example Request:
```bash
curl -X POST http://localhost:8080/v1/events \
-H "Content-Type: application/json" \
-d '{
  "user_id": "850e8400-e29b-41d4-a716-446655440000",
  "type": "click",
  "metadata": {"ip": "8.8.8.8"},
  "payload": {"button": "buy_now"}
}'
```

## 🔍 Monitoring & Verification
- **Logs**: View real-time processing logs.
```bash
docker logs -f sp_worker
```
- **Database**: Connect to PostgreSQL to see enriched data.
```bash
docker exec -it sp_postgres psql -U user -d streampulse -c "SELECT * FROM events;"
```
- **Queue**: Monitor message flow in the RabbitMQ Dashboard at http://localhost:15672.

## 🛡 Features Implemented
- **Asynchronous Processing**: Decoupled API and storage layers.
- **Data Enrichment**: Real-time GeoIP lookup for incoming events.
- **Idempotency**: Duplicate event detection using Redis.
- **Graceful Shutdown**: Safe handling of system signals to prevent data loss.
