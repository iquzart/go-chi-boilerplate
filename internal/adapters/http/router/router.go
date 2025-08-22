package router

import (
	"go-chi-boilerplate/internal/adapters/database/postgresql"
	custom "go-chi-boilerplate/internal/adapters/http/middleware"
	"go-chi-boilerplate/internal/adapters/http/routes"
	"log/slog"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func SetupRouter(serviceName string, logger *slog.Logger, db *postgresql.PostgresDB, jwtSecret string) *chi.Mux {
	r := chi.NewRouter()

	r.Use(otelhttp.NewMiddleware(serviceName))
	r.Use(custom.LoggingMiddleware(logger))
	r.Use(middleware.Recoverer)
	r.Use(custom.MetricsMiddleware)

	registerRoutes(r, db, logger, jwtSecret)
	return r
}

func registerRoutes(r chi.Router, db *postgresql.PostgresDB, logger *slog.Logger, jwtSecret string) {
	routes.AddSystemRoutes(r, db)
	routes.AddAPIRoutes(r, db, logger, jwtSecret)
}
