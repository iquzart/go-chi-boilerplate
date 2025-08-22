package user

import (
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
	// Hash password
	hashed, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user.Password = string(hashed)
	return u.repo.Create(user)
}

func (u *UserUsecase) GetUserByID(id string) (*entities.User, error) {
	return u.repo.GetByID(id)
}

func (u *UserUsecase) ListUsers(offset, limit int) ([]*entities.User, int, error) {
	return u.repo.List(offset, limit)
}

func (u *UserUsecase) UpdateUser(user *entities.User) (*entities.User, error) {
	return u.repo.Update(user)
}

func (u *UserUsecase) ChangeStatus(id string, status string) (*entities.User, error) {
	return u.repo.UpdateStatus(id, status)
}

func (u *UserUsecase) DeleteUser(id string) error {
	return u.repo.Delete(id)
}
