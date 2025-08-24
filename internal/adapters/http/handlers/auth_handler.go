package handlers

import (
	"encoding/json"
	"go-chi-boilerplate/internal/adapters/http/dto"
	"go-chi-boilerplate/internal/app/usecases/auth"
	"net/http"
)

// LoginHandler godoc
// @Summary      User login
// @Description  Login with email and password to get JWT token pair
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
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "invalid request"})
			return
		}

		// Get the context from the request
		ctx := r.Context()
		ip := r.RemoteAddr

		// Pass the context to the use case
		accessToken, refreshToken, err := authUC.Login(ctx, req.Email, req.Password, ip)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		resp := dto.LoginResponse{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(resp)
	}
}

// RefreshHandler godoc
// @Summary      Refresh access token
// @Description  Get a new access token using a valid refresh token
// @Tags         authentication
// @Accept       json
// @Produce      json
// @Param        refresh body object{user_id=string,refresh_token=string} true "Refresh token info"
// @Success      200 {object} map[string]string
// @Failure      401 {object} map[string]string
// @Router       /api/v1/auth/refresh [post]
func RefreshHandler(authUC *auth.AuthUsecase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			UserID       string `json:"user_id"`
			RefreshToken string `json:"refresh_token"`
		}
		json.NewDecoder(r.Body).Decode(&req)

		accessToken, err := authUC.RefreshAccessToken(r.Context(), req.UserID, req.RefreshToken)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"access_token": accessToken})
	}
}

// LogoutHandler godoc
// @Summary      Logout user
// @Description  Invalidate the user's refresh token (logout)
// @Tags         authentication
// @Accept       json
// @Produce      json
// @Param        logout body object{user_id=string} true "User ID"
// @Success      200 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /api/v1/auth/logout [post]
func LogoutHandler(authUC *auth.AuthUsecase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			UserID string `json:"user_id"`
		}
		json.NewDecoder(r.Body).Decode(&req)

		// Get the context from the request
		ctx := r.Context()

		ip := r.RemoteAddr

		// Pass the context to the use case.
		if err := authUC.Logout(ctx, req.UserID, ip); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "logged out"})
	}
}
