package main

import (
	"context"
	"go-chi-boilerplate/internal/adapters/database/postgresql"
	"go-chi-boilerplate/internal/adapters/http/server"
	"go-chi-boilerplate/internal/config"
	"go-chi-boilerplate/internal/meta"
	"log/slog"
	"time"
)

// Package docs contains the Swagger metadata for go-chi-boilerplate API.
//
// @title go-chi-boilerplate API
// @version 1.0.0
// @description This is the go-chi-boilerplate API documentation.
// @contact.name API Support
// @contact.url http://www.example.com/support
// @contact.email support@example.com
// @license.name MIT
// @license.url https://opensource.org/licenses/MIT
// @host localhost:8080
// @BasePath /
// @schemes http
//
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Enter the token with the `Bearer: ` prefix, e.g. "Bearer abcde12345".
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

	// Connect to PostgreSQL
	db, err := postgresql.New(cfg.Database, logger)
	if err != nil {
		meta.Fatal(logger, "failed to connect to database", "error", err)
	}
	defer db.Close()

	// Initialize PostgreSQL metrics
	meta.InitDBMetrics(db)

	// Optional: run migrations
	if err := postgresql.RunMigrations(db, "./migrations", logger); err != nil {
		meta.Fatal(logger, "failed to run migrations", "error", err)
	}

	// Start server
	server.New(cfg.Server, logger, db).Run()
}

func shutdownTracer(tp interface{ Shutdown(context.Context) error }, logger *slog.Logger) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := tp.Shutdown(ctx); err != nil {
		logger.Error("failed to shutdown tracer", "error", err)
	}
}
