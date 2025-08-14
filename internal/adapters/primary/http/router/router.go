package router

import (
	custom "go-chi-boilerplate/internal/adapters/primary/http/middleware"
	"go-chi-boilerplate/internal/adapters/primary/http/routes"
	"go-chi-boilerplate/internal/adapters/secondary/database/postgresql"
	"log/slog"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func SetupRouter(serviceName string, logger *slog.Logger, db *postgresql.PostgresDB) *chi.Mux {
	r := chi.NewRouter()

	r.Use(otelhttp.NewMiddleware(serviceName))
	r.Use(custom.LoggingMiddleware(logger))
	r.Use(middleware.Recoverer)
	r.Use(custom.MetricsMiddleware)

	registerRoutes(r, db)
	return r
}

func registerRoutes(r chi.Router, db *postgresql.PostgresDB) {
	routes.AddSystemRoutes(r, db)
	routes.AddApiRoutes(r)
}
