package repositories

import "go-chi-boilerplate/internal/core/entities"

type UserRepository interface {
	Create(user *entities.User) (*entities.User, error)
	GetByID(id string) (*entities.User, error)
	List(offset, limit int) ([]*entities.User, int, error)
	Update(user *entities.User) (*entities.User, error)
	UpdateStatus(id string, status string) (*entities.User, error)
	Delete(id string) error
	GetByEmail(email string) (*entities.User, error)
}
