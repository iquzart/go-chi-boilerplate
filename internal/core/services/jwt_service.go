package services

import (
	"time"

	"github.com/go-chi/jwtauth/v5"
)

type JWTService struct {
	TokenAuth *jwtauth.JWTAuth
}

// NewJWTService initializes the JWT service with the secret
func NewJWTService(secret string) *JWTService {
	return &JWTService{
		TokenAuth: jwtauth.New("HS256", []byte(secret), nil),
	}
}

// GenerateToken generates a JWT using user ID and role
func (s *JWTService) GenerateToken(userID string, role string) (string, error) {
	claims := map[string]interface{}{
		"sub":  userID,
		"role": role,
		"exp":  time.Now().Add(time.Hour).Unix(), // 1 hour expiry
	}

	_, tokenString, err := s.TokenAuth.Encode(claims)
	return tokenString, err
}
