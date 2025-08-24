package routes

import (
	"go-chi-boilerplate/internal/adapters/cache/redis"
	"go-chi-boilerplate/internal/adapters/database/postgresql"
	"go-chi-boilerplate/internal/adapters/http/handlers"
	"net/http"

	"github.com/go-chi/chi/v5"

	_ "go-chi-boilerplate/docs"
)

func AddSystemRoutes(r chi.Router, db *postgresql.PostgresDB, redisDB *redis.RedisDB) {
	system := chi.NewRouter()

	system.Get("/health", handlers.Health)
	system.Get("/liveness", handlers.Liveness)
	system.Get("/readiness", handlers.Readiness(db, redisDB))

	system.Handle("/metrics", handlers.MetricsHandler())

	system.Get("/swagger/openapi.json", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "docs/v3/openapi.json")
	})

	system.Handle("/swagger/*", handlers.SwaggerHandler())

	r.Mount("/system", system)
}
