package main

import (
	"context"
	"go-chi-boilerplate/internal/adapters/primary/http/server"
	"go-chi-boilerplate/internal/config"
	"go-chi-boilerplate/internal/meta"
	"time"
)

// Package docs contains the Swagger metadata for go-chi-boilerplate API.
// @title go-chi-boilerplate API
// @version 1.0
// @description This is the go-chi-boilerplate API documentation.
// @contact.name API Support
// @contact.url http://www.example.com/support
// @contact.email support@example.com
// @license.name MIT
// @license.url https://opensource.org/licenses/MIT
// @host localhost:8080
// @BasePath /
// @schemes http
func main() {
	cfg := config.GetServerConfigs()

	logger := meta.NewLogger(cfg.LogLevel)

	meta.InitMetrics()

	tp, err := meta.InitTracer(cfg.ServiceName, cfg.OTLPEndpoint, logger)
	if err != nil {
		logger.Error("failed to initialize tracer", "error", err)
	}
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := tp.Shutdown(ctx); err != nil {
			logger.Error("failed to shutdown tracer", "error", err)
		}
	}()

	server.New(cfg, logger).Run()
}
