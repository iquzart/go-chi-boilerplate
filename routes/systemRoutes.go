package routes

import (
	controllers "go-chi-boilerplate/controllers/system"

	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// addSystemRoutes adds the routes for the System controller to the specified router.
func addSystemRoutes(r chi.Router) {
	// Create a new sub-router for the System controller under the specified router.
	system := chi.NewRouter()

	// Add a GET route for the Health method of the System controller.
	system.Get("/health", controllers.Health)

	// Add a GET route for the Metrics endpoint using the Prometheus HTTP handler.
	system.Get("/metrics", func(w http.ResponseWriter, r *http.Request) {
		promhttp.Handler().ServeHTTP(w, r)
	})

	// Mount the sub-router to the specified router.
	r.Mount("/system", system)
}
