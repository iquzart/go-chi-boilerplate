package handlers

import (
	"context"
	"encoding/json"
	"go-chi-boilerplate/internal/adapters/secondary/database/postgresql"
	"net/http"
	"time"
)

// Health godoc
// @Summary Show the health status
// @Description Get the health status of the service
// @Tags system
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string
// @Router /system/health [get]
func Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]string{
		"status":  "ok",
		"message": "Working!",
	}

	json.NewEncoder(w).Encode(response)
}

// Liveness godoc
// @Summary Show if service is alive
// @Description Liveness probe for Kubernetes
// @Tags system
// @Accept json
// @Produce json
// @Success 200 {object} map[string]bool
// @Router /system/liveness [get]
func Liveness(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]bool{
		"alive": true,
	}

	json.NewEncoder(w).Encode(response)
}

// Readiness godoc
// @Summary Show if service is ready
// @Description Readiness probe for Kubernetes (checks DB connectivity)
// @Tags system
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 503 {object} map[string]interface{}
// @Router /system/readiness [get]
func Readiness(db *postgresql.PostgresDB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// Short timeout to avoid stressing DB
		ctx, cancel := context.WithTimeout(r.Context(), 1*time.Second)
		defer cancel()

		if err := db.DB.PingContext(ctx); err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"ready":     false,
				"db_status": "unreachable",
				"error":     err.Error(),
			})
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"ready":     true,
			"db_status": "ok",
		})
	}
}
