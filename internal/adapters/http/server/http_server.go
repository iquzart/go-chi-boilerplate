package server

import (
	"context"
	"go-chi-boilerplate/internal/adapters/cache/redis"
	"go-chi-boilerplate/internal/adapters/database/postgresql"
	"go-chi-boilerplate/internal/adapters/http/router"
	"go-chi-boilerplate/internal/config"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"
)

type Server struct {
	cfg    *config.ServerConfigs
	logger *slog.Logger
	http   *http.Server
}

// New creates a Server with all dependencies injected
func New(cfg *config.ServerConfigs, logger *slog.Logger, db *postgresql.PostgresDB, redisDB *redis.RedisDB) *Server {
	r := router.SetupRouter(cfg.ServiceName, logger, db, redisDB, cfg.JWTSecret)

	return &Server{
		cfg:    cfg,
		logger: logger,
		http: &http.Server{
			Addr:    cfg.Port,
			Handler: r,
		},
	}
}

// Run starts the HTTP server and handles graceful shutdown
func (s *Server) Run() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	go func() {
		s.logger.Info("service started", "port", s.cfg.Port)
		if err := s.http.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.Error("error starting server", "error", err)
		}
	}()

	<-stop
	s.logger.Info("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.http.Shutdown(ctx); err != nil {
		s.logger.Error("error during server shutdown", "error", err)
	}

	s.logger.Info("server stopped gracefully")
}
