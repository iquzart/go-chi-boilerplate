package auth

import (
	"context"
	"errors"
	"go-chi-boilerplate/internal/core/repositories"
	"go-chi-boilerplate/internal/core/services"
	"go-chi-boilerplate/internal/meta"
	"log/slog"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type AuthUsecase struct {
	UserRepo    repositories.UserRepository
	RefreshRepo repositories.RefreshTokenRepository
	JWT         *services.JWTService
	Logger      *slog.Logger
	// Add TTL to the struct to pass it to the repository
	TTL time.Duration
}

// NewAuthUsecase initializes AuthUsecase
func NewAuthUsecase(userRepo repositories.UserRepository, jwtService *services.JWTService, logger *slog.Logger, refreshRepo repositories.RefreshTokenRepository, ttl time.Duration) *AuthUsecase {
	return &AuthUsecase{
		UserRepo:    userRepo,
		RefreshRepo: refreshRepo,
		JWT:         jwtService,
		Logger:      logger,
		TTL:         ttl,
	}
}

func (a *AuthUsecase) Login(email, password, ip string) (string, string, error) {
	user, err := a.UserRepo.GetByEmail(email)
	if err != nil {
		meta.AuthEvent(a.Logger, "login_failed", "", email, "", ip, "invalid credentials")
		return "", "", errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		meta.AuthEvent(a.Logger, "login_failed", user.ID, email, user.Role, ip, "invalid credentials")
		return "", "", errors.New("invalid credentials")
	}

	accessToken, err := a.JWT.GenerateAccessToken(user.ID, user.Role)
	if err != nil {
		meta.AuthEvent(a.Logger, "login_failed", user.ID, email, user.Role, ip, "token generation failed")
		return "", "", errors.New("failed to generate access token")
	}

	refreshToken, err := a.JWT.GenerateRefreshToken(user.ID)
	if err != nil {
		meta.AuthEvent(a.Logger, "login_failed", user.ID, email, user.Role, ip, "refresh token generation failed")
		return "", "", errors.New("failed to generate refresh token")
	}

	// Corrected call to SaveRefreshToken
	if err := a.RefreshRepo.SaveRefreshToken(context.Background(), user.ID, refreshToken, time.Now().Add(a.TTL)); err != nil {
		meta.AuthEvent(a.Logger, "login_failed", user.ID, email, user.Role, ip, "saving refresh token failed")
		return "", "", errors.New("failed to save refresh token")
	}

	meta.AuthEvent(a.Logger, "login_success", user.ID, email, user.Role, ip, "")
	return accessToken, refreshToken, nil
}

// RefreshAccessToken handles refreshing the access token
func (a *AuthUsecase) RefreshAccessToken(ctx context.Context, userID, token string) (string, error) {
	storedToken, err := a.RefreshRepo.GetRefreshToken(ctx, userID)
	if err != nil || storedToken != token {
		return "", errors.New("invalid or expired refresh token")
	}

	user, err := a.UserRepo.GetByID(userID)
	if err != nil {
		return "", errors.New("user not found")
	}

	// Corrected call to GenerateAccessToken
	newAccessToken, err := a.JWT.GenerateAccessToken(user.ID, user.Role)
	if err != nil {
		return "", errors.New("failed to generate new access token")
	}

	return newAccessToken, nil
}

// Logout removes the refresh token from the repository
func (a *AuthUsecase) Logout(ctx context.Context, userID string) error {
	return a.RefreshRepo.DeleteRefreshToken(ctx, userID)
}
