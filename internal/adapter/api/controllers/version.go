package controllers

import (
	"encoding/json"
	"net/http"
	"os"
)

// APIVersion returns the current API version
// @Summary Get API version
// @Description Returns the current version of the API in JSON format
// @Tags api
// @Produce json
// @Success 200 {object} map[string]string "API version"
// @Router /api/version [get]
func APIVersion(w http.ResponseWriter, r *http.Request) {
	apiVersion := os.Getenv("API_VERSION")
	if apiVersion == "" {
		apiVersion = "v1.0.0"
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]string{
		"version": apiVersion,
	}

	json.NewEncoder(w).Encode(response)
}
