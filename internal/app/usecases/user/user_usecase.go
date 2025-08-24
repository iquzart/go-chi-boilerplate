package user

import (
	"context"
	"errors"
	"go-chi-boilerplate/internal/core/entities"
	"go-chi-boilerplate/internal/core/repositories"

	"golang.org/x/crypto/bcrypt"
)

type UserUsecase struct {
	repo repositories.UserRepository
}

func NewUserUsecase(repo repositories.UserRepository) *UserUsecase {
	return &UserUsecase{repo: repo}
}

func (u *UserUsecase) CreateUser(ctx context.Context, user *entities.User) (*entities.User, error) {
	// Check if email exists
	exists, err := u.repo.ExistsByEmail(ctx, user.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("email already exists")
	}

	// Hash password
	hashed, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user.Password = string(hashed)
	return u.repo.Create(ctx, user)
}

func (u *UserUsecase) UpdateUser(ctx context.Context, user *entities.User) (*entities.User, error) {
	exists, err := u.repo.ExistsByEmail(ctx, user.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		// If the existing user has a different ID, prevent update
		existingUser, _ := u.repo.GetByEmail(ctx, user.Email)
		if existingUser.ID != user.ID {
			return nil, errors.New("email already exists")
		}
	}

	return u.repo.Update(ctx, user)
}

func (u *UserUsecase) GetUserByID(ctx context.Context, id string) (*entities.User, error) {
	return u.repo.GetByID(ctx, id)
}

func (u *UserUsecase) ListUsers(ctx context.Context, offset, limit int) ([]*entities.User, int, error) {
	return u.repo.List(ctx, offset, limit)
}

func (u *UserUsecase) ChangeStatus(ctx context.Context, id string, status string) (*entities.User, error) {
	return u.repo.UpdateStatus(ctx, id, status)
}

func (u *UserUsecase) DeleteUser(ctx context.Context, id string) error {
	return u.repo.Delete(ctx, id)
}
