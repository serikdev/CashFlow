

# 💸 CashFlow

![Go](https://img.shields.io/badge/Go-1.24+-00ADD8?logo=go\&logoColor=white)
![Postgres](https://img.shields.io/badge/PostgreSQL-15-4169E1?logo=postgresql\&logoColor=white)
![Kafka](https://img.shields.io/badge/Apache%20Kafka-Event%20Streaming-231F20?logo=apachekafka\&logoColor=white)
![Docker](https://img.shields.io/badge/Docker-Ready-2496ED?logo=docker\&logoColor=white)
![CI/CD](https://img.shields.io/badge/GitHub-Actions-2088FF?logo=githubactions\&logoColor=white)

**CashFlow** is a financial transaction management system built with **Clean Architecture**,
using **Go + PostgreSQL + Kafka**.

It provides APIs for managing accounts and transactions (**deposit, withdraw, transfer**),
with **event-driven processing** via Kafka consumers and producers.

---

## 🚀 Features

* ✅ **Clean Architecture** (Handlers → Usecases → Repositories → Adapters)
* ✅ **PostgreSQL + pgxpool** for persistence
* ✅ **Kafka producers/consumers** for asynchronous transaction processing
* ✅ **REST API with Swagger docs**
* ✅ **Graceful shutdown** with context & signals
* ✅ **Goose migrations** for DB schema versioning

---

## 🏗️ Project Architecture

```text
cmd/            – Application entrypoint
deployment/     - Deployment
internal/
  adapter/      – DB Repositories (Postgres with pgxpool)
  entity/       – Domain models
  kafka/        – Kafka Producer & Consumer
  port/rest/    – HTTP Handlers (Swagger-ready)
  usecase/      – Business logic (Account & Transaction services)
pkg/
  database/     – DB pool initialization
  logger/       – Logrus-based structured logging
docs/           – Auto-generated Swagger docs
```

---

## ⚙️ Tech Stack

* **Language**: Go (1.24+)
* **Database**: PostgreSQL 15
* **Message Broker**: Apache Kafka
* **Migrations**: Goose
* **Logging**: Logrus
* **API Docs**: Swaggo + Swagger UI

---

## 📦 Installation

### 1. Clone repository

```bash
git clone https://github.com/serikdev/CashFlow.git
cd CashFlow
```

### 2. Run PostgreSQL + Kafka with Docker

```bash
docker-compose up -d
```

### 3. Run database migrations

```bash
goose -dir migrations postgres "postgres://user:password@localhost:5432/cashflow?sslmode=disable" up
```

### 4. Run the application

```bash
go run cmd/api/main.go
```

---

## 📖 API Documentation

After starting the server, open Swagger UI:

👉 [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)

Example endpoints:

* `POST /api/accounts` → Create account
* `GET /api/accounts/{id}` → Get account
* `POST /api/accounts/{id}/deposit` → Deposit money
* `POST /api/accounts/{id}/withdraw` → Withdraw money
* `POST /api/accounts/{id}/transfer` → Transfer money
* `GET /api/accounts/{id}/transactions` → Transaction history

---

## 🔄 Kafka Integration

The system uses **Kafka topics**:

* `account-deposit` → Deposit events
* `account-withdraw` → Withdrawal events
* `account-transfer` → Transfer events

Producers publish transaction events,
Consumers subscribe and update database state accordingly.

---

## 🛠️ Example Requests

### Create Account

```bash
curl -X POST http://localhost:8080/api/accounts \
-H "Content-Type: application/json" \
-d '{
  "balance": 1000.50,
  "currency": "TMT"
}'
```

### Deposit

```bash
curl -X POST http://localhost:8080/api/accounts/1/deposit \
-H "Content-Type: application/json" \
-d '{"amount": 500}'
```

---

## 🧪 Running Tests

```bash
go test ./...
```

---

## 🧩 Roadmap

* [ ] Dockerfile & Helm Charts for Kubernetes
* [ ] gRPC API
* [ ] Redis cache for faster reads
* [ ] CI/CD GitHub Actions pipeline
* [ ] Monitoring (Prometheus + Grafana)

---

## 📜 License

MIT License © 2025 \[Serdar]
