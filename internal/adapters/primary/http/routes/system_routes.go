package routes

import (
	"go-chi-boilerplate/internal/adapters/primary/http/handlers"

	"github.com/go-chi/chi/v5"

	_ "go-chi-boilerplate/docs"
)

func AddSystemRoutes(r chi.Router) {
	system := chi.NewRouter()

	system.Get("/health", handlers.Health)

	system.Handle("/metrics", handlers.MetricsHandler())

	system.Handle("/swagger/*", handlers.SwaggerHandler())

	r.Mount("/system", system)
}
