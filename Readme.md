# RideMarket – Clean Architecture Backend in Go

Modern backend for a ride-sharing / ride-hailing marketplace (inspired by Snapp, Tapsi, Uber-like systems) built with **Go** following **Clean Architecture** principles.

Goal: separation of concerns, high testability, framework independence, clear domain logic.



## ✨ Main Features (current / planned)

- Passenger & Driver authentication (JWT + refresh token)
- User profile management
- Real-time driver location update (simulation or via websocket)
- Ride creation & request
- Nearby driver search & matching logic
- Ride lifecycle: requested → found → accepted → started → completed → canceled
- Basic fare estimation
- Ride history
- Rating & review after ride

## 🛠 Tech Stack

| Category            | Technology / Library                     |
|---------------------|------------------------------------------|
| Language            | Go 1.21+                                 |
| HTTP framework      | Gin / Echo / Chi / standard net/http     |
| Database            | PostgreSQL                               |
| ORM / Query         | GORM / sqlc / bun                        |
| Validation          | go-playground/validator.v10              |
| Authentication      | JWT (github.com/golang-jwt/jwt/v5)       |
| Config              | Viper / env / koanf                      |
| Logger              | zerolog / zap                            |
| Testing             | testing + testify / httptest             |
| Container           | Docker + docker-compose                  |

## 🚀 Quick Start

### Prerequisites

- Go 1.21+
- PostgreSQL 15+
- Redis
- Elk(Kibana&ElasticSearch&Filebeat)
- Prometheus & Grafana
- (optional) Docker & docker-compose

### Local development

```bash
# 1. Clone project
git clone https://github.com/heydarabadi/RideMarket-CleanWebApi-GoLang.git
cd RideMarket-CleanWebApi-GoLang

# 2. Copy environment file
cp .env.example .env

# 3. Adjust .env (database credentials, jwt secret, ...)

# 4. Download dependencies
go mod tidy

# 5. Run migrations (if using goose / golang-migrate / gorm auto-migrate)
# Example with goose:
# goose -dir infrastructure/persistence/migrations postgres "your-db-url" up

# 6. Start server
go run ./Src/main.go
# or with hot-reload
air 
Server usually starts at:
http://localhost:8080


Using Docker (recommended)
docker-compose up -d --build
```

