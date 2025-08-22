package auth

import (
	"errors"
	"go-chi-boilerplate/internal/core/repositories"
	"go-chi-boilerplate/internal/core/services"

	"golang.org/x/crypto/bcrypt"
)

type AuthUsecase struct {
	UserRepo repositories.UserRepository
	JWT      *services.JWTService
}

func NewAuthUsecase(repo repositories.UserRepository, jwt *services.JWTService) *AuthUsecase {
	return &AuthUsecase{
		UserRepo: repo,
		JWT:      jwt,
	}
}

func (a *AuthUsecase) Login(email, password string) (string, error) {
	user, err := a.UserRepo.GetByEmail(email)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", errors.New("invalid credentials")
	}

	token, err := a.JWT.GenerateToken(user.ID, user.Role)
	if err != nil {
		return "", errors.New("failed to generate token")
	}

	return token, nil
}
