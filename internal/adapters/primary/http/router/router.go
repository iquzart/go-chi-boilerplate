package router

import (
	custom "go-chi-boilerplate/internal/adapters/primary/http/middleware"
	"go-chi-boilerplate/internal/adapters/primary/http/routes"
	"log/slog"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func SetupRouter(serviceName string, logger *slog.Logger) *chi.Mux {
	r := chi.NewRouter()

	r.Use(otelhttp.NewMiddleware(serviceName))
	r.Use(custom.LoggingMiddleware(logger))
	r.Use(middleware.Recoverer)
	r.Use(custom.MetricsMiddleware)

	registerRoutes(r)
	return r
}

func registerRoutes(r chi.Router) {
	routes.AddSystemRoutes(r)
	routes.AddApiRoutes(r)
}
