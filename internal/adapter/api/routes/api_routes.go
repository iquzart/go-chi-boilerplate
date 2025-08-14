package routes

import (
	controllers "go-chi-boilerplate/internal/adapter/api/controllers"

	"github.com/go-chi/chi/v5"
)

func AddApiRoutes(rg chi.Router) {
	api := chi.NewRouter()

	api.Get("/version", controllers.APIVersion)

	rg.Mount("/api", api)
}
