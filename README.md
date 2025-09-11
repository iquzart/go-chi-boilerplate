# Go Chi Boilerplate

[![Build and Test](https://github.com/iquzart/go-chi-boilerplate/actions/workflows/ci.yaml/badge.svg?branch=main)](https://github.com/iquzart/go-chi-boilerplate/actions/workflows/ci.yaml)
[![Go Report Card](https://goreportcard.com/badge/github.com/iquzart/go-chi-boilerplate)](https://goreportcard.com/report/github.com/iquzart/go-chi-boilerplate)
![GitHub](https://img.shields.io/github/license/iquzart/go-chi-boilerplate)
![Metrics Support](https://img.shields.io/badge/Metrics%20Support-Prometheus-blue)

Golang Chi boilerplate for building microservices with **structured JSON logging**, **Prometheus metrics**, **health checks**, **Swagger documentation**, and **OpenTelemetry tracing**.

---

## Table of Contents

- [Features](#features)
- [Environment Variables](#environment-variables)
- [Running the Service](#running-the-service)
- [Docker Compose Services](#docker-compose-services)
- [Makefile Commands](#makefile-commands)
- [Endpoints](#endpoints)
- [License](#license)

---

## Features

1. Configurable server ports
2. Kubernetes-friendly health check endpoint
3. Prometheus metrics for monitoring
4. Swagger API documentation
5. OpenTelemetry distributed tracing
6. Custom structured JSON logging using `slog`
7. Ready-to-use project structure following **hexagonal architecture**

---

## Environment Variables

| Variable         | Description                                        | Default         |
|-----------------|----------------------------------------------------|----------------|
| `SERVICE_NAME`   | Name of the service used for tracing and logging   | `go-chi-boilerplate`   |
| `LOG_LEVEL`      | Logging level (`debug`, `info`, `warn`, `error`) | `info`         |
| `OTLP_ENDPOINT`  | OpenTelemetry collector endpoint                  | `otel-collector:4317` |
| `API_VERSION`    | API version returned by `/system/version`        | `v1.0.0`       |
| `PORT`           | Port on which the server listens                  | `8080`         |
| `DB_HOST`        | PostgreSQL host                                   | `postgres`     |
| `DB_PORT`        | PostgreSQL port                                   | `5432`         |
| `DB_USER`        | PostgreSQL user                                   | `appuser`      |
| `DB_PASSWORD`    | PostgreSQL password                               | `apppassword`  |
| `DB_NAME`        | PostgreSQL database name                           | `appdb`        |
| `DB_SSLMODE`     | PostgreSQL SSL mode                                | `disable`      |
| `REDIS_HOST`     | Redis host                                        | `redis`        |
| `REDIS_PORT`     | Redis port                                        | `6379`         |
| `REDIS_PASSWORD` | Redis password                                    | `redis!password123` |
| `REDIS_DB`       | Redis database number                             | `0`            |
| `REDIS_PREFIX`   | Redis key prefix                                  | `go-chi-boilerplate` |
| `REDIS_TTL`      | Default TTL for Redis keys (supports `15m`, `1h`) | `1h`           |

---

## Running the Service

### Local

```bash
git clone https://github.com/iquzart/go-chi-boilerplate.git
cd go-chi-boilerplate
make run
````

### Docker

```bash
make docker-build
make docker-run
```

### Full Stack with Docker Compose

```bash
make compose-up
# Access services:
# API: http://localhost:8080
# Grafana: http://localhost:3000
# Prometheus: http://localhost:9090
# Tempo: http://localhost:3200
```

Stop services:

```bash
make compose-down
```

---

## Docker Compose Services

- `go-chi-boilerplate`: The main Go service
- `postgres`: PostgreSQL database
- `redis`: Redis cache
- `tempo`: Grafana Tempo for tracing
- `otel-collector`: OpenTelemetry Collector
- `prometheus`: Metrics server
- `grafana`: Visualization and dashboards

Volumes are used for persistent storage:

- `postgres_data`
- `redis_data`
- `tempo_data`

---

## Makefile Commands

| Command             | Description                                          |
| ------------------- | ---------------------------------------------------- |
| `make help`         | Show this help message                               |
| `make clean`        | Remove binary and Docker resources                   |
| `make build`        | Build Go binary                                      |
| `make run`          | Run the binary locally (also generates Swagger docs) |
| `make docker-build` | Build Docker image                                   |
| `make docker-run`   | Run Docker container                                 |
| `make docker-stop`  | Stop Docker container                                |
| `make compose-up`   | Start all services via Docker Compose                |
| `make compose-down` | Stop all services via Docker Compose                 |
| `make swagger`      | Generate Swagger 2.0 and OpenAPI v3 docs             |

---

## Endpoints

| Endpoint          | Description                 |
| ----------------- | --------------------------- |
| `/system/health`  | Health check for Kubernetes |
| `/system/metrics` | Prometheus metrics          |
| `/api/version` | Returns API version         |
| `/swagger/*`      | Swagger documentation       |

---

## License

MIT

## Author Information

Muhammed Iqbal <iquzart@hotmail.com>
