package repositories

import "go-chi-boilerplate/internal/core/entities"

// UserRepository defines the interface for user persistence operations
type UserRepository interface {
	// Create inserts a new user into the repository
	Create(user *entities.User) (*entities.User, error)

	// GetByID retrieves a user by their unique ID
	GetByID(id string) (*entities.User, error)

	// List returns a list of users with pagination support
	// offset: starting index, limit: maximum number of users to return
	// Returns the list of users and the total number of users
	List(offset, limit int) ([]*entities.User, int, error)

	// Update modifies an existing user's information
	Update(user *entities.User) (*entities.User, error)

	// UpdateStatus changes the status of a user (e.g., active/inactive)
	UpdateStatus(id string, status string) (*entities.User, error)

	// Delete removes a user from the repository by ID
	Delete(id string) error

	// GetByEmail retrieves a user by their unique email
	GetByEmail(email string) (*entities.User, error)

	// ExistsByEmail checks if a user with the given email already exists
	// Returns true if the email exists, false otherwise
	ExistsByEmail(email string) (bool, error)
}
