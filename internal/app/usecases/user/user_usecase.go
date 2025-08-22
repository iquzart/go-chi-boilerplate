package user

import (
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

func (u *UserUsecase) CreateUser(user *entities.User) (*entities.User, error) {
	// Check if email exists
	exists, err := u.repo.ExistsByEmail(user.Email)
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
	return u.repo.Create(user)
}

func (u *UserUsecase) UpdateUser(user *entities.User) (*entities.User, error) {
	exists, err := u.repo.ExistsByEmail(user.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		// If the existing user has a different ID, prevent update
		existingUser, _ := u.repo.GetByEmail(user.Email)
		if existingUser.ID != user.ID {
			return nil, errors.New("email already exists")
		}
	}

	return u.repo.Update(user)
}

func (u *UserUsecase) GetUserByID(id string) (*entities.User, error) {
	return u.repo.GetByID(id)
}

func (u *UserUsecase) ListUsers(offset, limit int) ([]*entities.User, int, error) {
	return u.repo.List(offset, limit)
}

func (u *UserUsecase) ChangeStatus(id string, status string) (*entities.User, error) {
	return u.repo.UpdateStatus(id, status)
}

func (u *UserUsecase) DeleteUser(id string) error {
	return u.repo.Delete(id)
}
