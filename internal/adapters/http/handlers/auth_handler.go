package handlers

import (
	"encoding/json"
	"go-chi-boilerplate/internal/adapters/http/dto"
	"go-chi-boilerplate/internal/app/usecases/auth"
	"net/http"
)

// LoginHandler godoc
// @Summary      User login
// @Description  Login with email and password to get JWT token
// @Tags         authentication
// @Accept       json
// @Produce      json
// @Param        login body dto.LoginRequest true "Login info"
// @Success      200 {object} dto.LoginResponse
// @Failure      400 {object} map[string]string
// @Failure      401 {object} map[string]string
// @Router       /api/v1/auth/login [post]
func LoginHandler(authUC *auth.AuthUsecase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req dto.LoginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "invalid request"})
			return
		}

		// Capture client IP
		ip := r.RemoteAddr

		token, err := authUC.Login(req.Email, req.Password, ip)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		resp := dto.LoginResponse{AccessToken: token}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
	}
}
