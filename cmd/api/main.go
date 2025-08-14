package main

import (
	"context"
	"go-chi-boilerplate/internal/adapters/primary/http/server"
	"go-chi-boilerplate/internal/config"
	"go-chi-boilerplate/internal/meta"
	"log/slog"
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
	// Load configs
	cfg, err := config.GetAppConfigs()
	if err != nil {
		meta.Fatal(meta.NewLogger("error"), "failed to load application configs", "error", err)
	}

	// Init logger
	logger := meta.NewLogger(cfg.Server.LogLevel)

	// Init metrics
	meta.InitMetrics()

	// Init tracer
	tp, err := meta.InitTracer(cfg.Server.ServiceName, cfg.Server.OTLPEndpoint, logger)
	if err != nil {
		meta.Fatal(logger, "failed to initialize tracer", "error", err)
	}
	defer shutdownTracer(tp, logger)

	// Start server
	server.New(cfg.Server, logger).Run()
}

func shutdownTracer(tp interface{ Shutdown(context.Context) error }, logger *slog.Logger) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := tp.Shutdown(ctx); err != nil {
		logger.Error("failed to shutdown tracer", "error", err)
	}
}
