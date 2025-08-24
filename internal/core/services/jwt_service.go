package services

import (
	"time"

	"github.com/go-chi/jwtauth/v5"
	"github.com/google/uuid"
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

// GenerateAccessToken generates a JWT using user ID and role with a short expiry
func (s *JWTService) GenerateAccessToken(userID string, role string) (string, error) {
	claims := map[string]interface{}{
		"sub":  userID,
		"role": role,
		"exp":  time.Now().Add(3 * time.Minute).Unix(), // 3-minute expiry
	}

	_, tokenString, err := s.TokenAuth.Encode(claims)
	return tokenString, err
}

// GenerateRefreshToken generates a simple JWT without an explicit expiration claim.
// The TTL is managed by the Redis store.
func (s *JWTService) GenerateRefreshToken(userID string) (string, error) {
	// A unique ID (JTI) is crucial for token invalidation in Redis.
	jti := uuid.NewString()

	claims := map[string]interface{}{
		"sub": userID,
		"jti": jti,
	}

	_, tokenString, err := s.TokenAuth.Encode(claims)
	return tokenString, err
}
