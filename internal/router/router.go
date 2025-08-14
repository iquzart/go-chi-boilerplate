package router

import (
	"go-chi-boilerplate/internal/adapter/api/routes"
	"go-chi-boilerplate/internal/meta"
	"log/slog"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func SetupRouter(serviceName string, logger *slog.Logger) *chi.Mux {
	r := chi.NewRouter()

	r.Use(otelhttp.NewMiddleware(serviceName))
	r.Use(meta.LoggingMiddleware(logger))
	r.Use(meta.MetricsMiddleware)
	r.Use(middleware.Recoverer)

	registerRoutes(r)
	return r
}

func registerRoutes(r chi.Router) {
	routes.AddSystemRoutes(r)
	routes.AddApiRoutes(r)
}
