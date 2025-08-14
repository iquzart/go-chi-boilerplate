package controllers

import (
	"encoding/json"
	"net/http"
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
