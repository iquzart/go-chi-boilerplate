package controllers

import (
	"net/http"
)

func Health(w http.ResponseWriter, r *http.Request) {
	// Respond with a 200 status code and "Working!" message.
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Working!"))
}
