package handlers

import (
	"encoding/json"
	"fmt"
	"go-chi-boilerplate/internal/adapters/http/dto"
	"go-chi-boilerplate/internal/app/usecases/user"
	"net/http"

	"github.com/go-chi/jwtauth/v5"
)

// GetProfile godoc
// @Summary      Get current user profile
// @Description  Returns the profile of the currently authenticated user
// @Tags         userProfile
// @Security     BearerAuth
// @Produce      json
// @Success      200 {object} dto.ProfileResponse
// @Failure      401 {object} map[string]string
// @Router       /api/v1/profile [get]
func GetProfile(userUC *user.UserUsecase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, claims, _ := jwtauth.FromContext(r.Context())

		// temp
		claimsJSON, _ := json.Marshal(claims)
		fmt.Println("JWT Claims:", string(claimsJSON))

		userID, ok := claims["sub"].(string)
		if !ok || userID == "" {
			http.Error(w, `{"error":"invalid token claims"}`, http.StatusUnauthorized)
			return
		}

		// Pass the context from the request
		user, err := userUC.GetUserByID(r.Context(), userID)
		if err != nil {
			http.Error(w, `{"error":"user not found"}`, http.StatusNotFound)
			return
		}

		resp := dto.ProfileResponse{
			UserID:    user.ID,
			Role:      user.Role,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Email:     user.Email,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
	}
}
