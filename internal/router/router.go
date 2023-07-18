package router

import (
	routes "go-chi-boilerplate/internal/interfaces/api/routes"
	customMiddleware "go-chi-boilerplate/internal/middlewares"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// InitRouter initializes the Chi router and sets up the middleware and routes.
func SetupRouter(serviceName string) *chi.Mux {

	// Create a new Chi router instance.
	router := chi.NewRouter()

	// Use the Chi logger middleware to log HTTP requests and responses.
	router.Use(middleware.Logger)

	// Use the Chi recovery middleware to recover from panics and return an error response instead of crashing.
	router.Use(middleware.Recoverer)

	// Use the Prometheus middleware to instrument the router with metrics.
	router.Use(customMiddleware.PrometheusMetrics(serviceName))

	// Add the routes to router.
	addRoutes(router)

	// Return the initialized Chi router.
	return router
}

// getRoutes adds the system and application routes to the specified router.
func addRoutes(r chi.Router) {
	// Add the routes for the System.
	routes.AddSystemRoutes(r)
	// Add the routes for app.
	routes.AddAPIRoutes(r)
}
