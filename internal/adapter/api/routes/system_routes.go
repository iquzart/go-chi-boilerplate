package routes

import (
	controllers "go-chi-boilerplate/internal/adapter/api/controllers"
	"go-chi-boilerplate/internal/meta"

	"github.com/go-chi/chi/v5"

	_ "go-chi-boilerplate/docs"
)

func AddSystemRoutes(r chi.Router) {
	system := chi.NewRouter()

	system.Get("/health", controllers.Health)

	system.Handle("/metrics", meta.MetricsHandler())

	system.Handle("/swagger/*", controllers.SwaggerHandler())

	r.Mount("/system", system)
}
