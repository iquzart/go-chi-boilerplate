package routes

import (
	"go-chi-boilerplate/internal/adapters/primary/http/handlers"
	"go-chi-boilerplate/internal/adapters/secondary/database/postgresql"

	"github.com/go-chi/chi/v5"

	_ "go-chi-boilerplate/docs"
)

func AddSystemRoutes(r chi.Router, db *postgresql.PostgresDB) {
	system := chi.NewRouter()

	system.Get("/health", handlers.Health)
	system.Get("/liveness", handlers.Liveness)
	system.Get("/readiness", handlers.Readiness(db))

	system.Handle("/metrics", handlers.MetricsHandler())

	system.Handle("/swagger/*", handlers.SwaggerHandler())

	r.Mount("/system", system)
}
