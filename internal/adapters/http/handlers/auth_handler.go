package handlers

import (
	"go-chi-boilerplate/internal/adapters/http/dto"
	"go-chi-boilerplate/internal/app/usecases/auth"
	"encoding/json"
	"net/http"
)

// LoginHandler godoc
// @Summary      User login
// @Description  Login with email and password to get JWT token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        login body dto.LoginRequest true "Login info"
// @Success      200 {object} dto.LoginResponse
// @Failure      400 {object} map[string]string
// @Failure      401 {object} map[string]string
// @Router       /api/v1/auth/login [post]
func LoginHandler(authUsecase *auth.AuthUsecase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req dto.LoginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		token, err := authUsecase.Login(req.Email, req.Password)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		resp := dto.LoginResponse{AccessToken: token}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}
