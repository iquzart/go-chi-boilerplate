# Go Chi Boilerplate

[![Build and Test](https://github.com/iquzart/go-chi-boilerplate/actions/workflows/ci.yaml/badge.svg?branch=main)](https://github.com/iquzart/go-chi-boilerplate/actions/workflows/ci.yaml)
[![Go Report Card](https://goreportcard.com/badge/github.com/iquzart/go-chi-boilerplate)](https://goreportcard.com/report/github.com/iquzart/go-chi-boilerplate)
![GitHub](https://img.shields.io/github/license/iquzart/go-chi-boilerplate)
![Metrics Support](https://img.shields.io/badge/Metrics%20Support-Prometheus-blue)

Golang Chi Boilerplate for building microservices with structured logging, Prometheus metrics, health checks, and tracing.

## Features

1. Custom ports
2. Health check endpoint for Kubernetes
3. Prometheus metrics
4. Swagger documentation
5. OpenTelemetry tracing
6. Custom structured JSON logging using `slog`

## Environment Variables

| Variable                  | Description                                       | Default   |
|---------------------------|---------------------------------------------------|----------|
| `SERVICE_NAME`            | Name of the service used for tracing and logging  | `go-chi-app` |
| `LOG_LEVEL`               | Logging level (debug, info, warn, error)         | `info`    |
| `OTLP_ENDPOINT`           | OpenTelemetry collector endpoint                 | `localhost:4317` |
| `API_VERSION`             | API version returned by `/system/version`        | `v1.0.0` |
| `PORT`                    | Port on which the server listens                 | `8080`    |

## License

MIT

## Author Information

Muhammed Iqbal <iquzart@hotmail.com>

