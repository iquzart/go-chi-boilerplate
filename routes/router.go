package routes

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

// InitRouter initializes the Chi router and sets up the middleware and routes.
func InitRouter() *chi.Mux {
	// Create a new Chi router instance.
	router := chi.NewRouter()

	// Use the Chi logger middleware to log HTTP requests and responses.
	router.Use(middleware.Logger)

	// Use the Chi recovery middleware to recover from panics and return an error response instead of crashing.
	router.Use(middleware.Recoverer)

	// Add the routes to router.
	getRoutes(router)

	// Return the initialized Chi router.
	return router
}

// getRoutes adds the system and application routes to the specified router.
func getRoutes(r chi.Router) {
	// Add the routes for the System.
	addSystemRoutes(r)
	// Add the routes for app.
	addAPIRoutes(r)
}