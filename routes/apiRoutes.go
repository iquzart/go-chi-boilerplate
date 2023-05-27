package routes

import (
	controllers "go-chi-boilerplate/controllers/api"

	"github.com/go-chi/chi/v5"
)

// addAPIRoutes adds the routes for the API controller to the specified router group.
func addAPIRoutes(rg chi.Router) {
	// Create a new group of routes for the API controller under the specified router group.
	api := chi.NewRouter()

	// Add a GET route for the APIVersion method of the API controller.
	api.Get("/version", controllers.APIVersion)

	// Mount the API routes under the specified router group.
	rg.Mount("/api", api)
}
