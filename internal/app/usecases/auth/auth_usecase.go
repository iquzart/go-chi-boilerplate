package auth

import (
	"errors"
	"go-chi-boilerplate/internal/core/repositories"
	"go-chi-boilerplate/internal/core/services"
	"go-chi-boilerplate/internal/meta"
	"log/slog"

	"golang.org/x/crypto/bcrypt"
)

// AuthUsecase handles authentication logic
type AuthUsecase struct {
	UserRepo repositories.UserRepository
	JWT      *services.JWTService
	Logger   *slog.Logger
}

// NewAuthUsecase creates a new AuthUsecase instance
func NewAuthUsecase(repo repositories.UserRepository, jwt *services.JWTService, logger *slog.Logger) *AuthUsecase {
	return &AuthUsecase{
		UserRepo: repo,
		JWT:      jwt,
		Logger:   logger,
	}
}

// Login authenticates a user and generates a JWT token.
// Logs authentication events for auditing.
func (a *AuthUsecase) Login(email, password, ip string) (string, error) {
	user, err := a.UserRepo.GetByEmail(email)
	if err != nil {
		meta.AuthEvent(a.Logger, "login_failed", "", email, "", ip, "invalid credentials")
		return "", errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		meta.AuthEvent(a.Logger, "login_failed", user.ID, email, user.Role, ip, "invalid credentials")
		return "", errors.New("invalid credentials")
	}

	token, err := a.JWT.GenerateToken(user.ID, user.Role)
	if err != nil {
		meta.AuthEvent(a.Logger, "login_failed", user.ID, email, user.Role, ip, "token generation failed")
		return "", errors.New("failed to generate token")
	}

	meta.AuthEvent(a.Logger, "login_success", user.ID, email, user.Role, ip, "")
	return token, nil
}
