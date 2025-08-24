package router

import (
	"go-chi-boilerplate/internal/adapters/cache/redis"
	"go-chi-boilerplate/internal/adapters/database/postgresql"
	custom "go-chi-boilerplate/internal/adapters/http/middleware"
	"go-chi-boilerplate/internal/adapters/http/routes"
	"go-chi-boilerplate/internal/config"
	"log/slog"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func SetupRouter(cfg *config.AppConfigs, logger *slog.Logger, db *postgresql.PostgresDB, redisDB *redis.RedisDB) *chi.Mux {
	r := chi.NewRouter()

	r.Use(otelhttp.NewMiddleware(cfg.Server.ServiceName))
	r.Use(custom.LoggingMiddleware(logger))
	r.Use(middleware.Recoverer)
	r.Use(custom.MetricsMiddleware)

	registerRoutes(r, db, logger, redisDB, cfg)
	return r
}

func registerRoutes(r chi.Router, db *postgresql.PostgresDB, logger *slog.Logger, redisDB *redis.RedisDB, cfg *config.AppConfigs) {
	routes.AddSystemRoutes(r, db, redisDB)
	routes.AddAPIRoutes(r, db, logger, redisDB, cfg)
}
