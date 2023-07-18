package controllers

import (
	"net/http"
	"os"
)

func APIVersion(w http.ResponseWriter, r *http.Request) {
	// Get the API version from the "API_VERSION" environment variable, or use a default value of "v1.0.0" if not set.
	apiVersion := os.Getenv("API_VERSION")
	if apiVersion == "" {
		apiVersion = "v1.0.0"
	}

	// Respond with a 200 status code and the API version as the response body.
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(apiVersion))
}
