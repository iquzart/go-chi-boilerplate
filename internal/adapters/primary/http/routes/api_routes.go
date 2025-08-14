package routes

import (
	"go-chi-boilerplate/internal/adapters/primary/http/handlers"

	"github.com/go-chi/chi/v5"
)

func AddApiRoutes(rg chi.Router) {
	api := chi.NewRouter()

	api.Get("/version", handlers.APIVersion)

	rg.Mount("/api", api)
}
