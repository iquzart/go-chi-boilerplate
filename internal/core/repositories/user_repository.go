package repositories

import (
	"context"
	"go-chi-boilerplate/internal/core/entities"
	"time"
)

// UserRepository defines the interface for user persistence operations
type UserRepository interface {
	// Create inserts a new user into the repository
	Create(ctx context.Context, user *entities.User) (*entities.User, error)

	// GetByID retrieves a user by their unique ID
	GetByID(ctx context.Context, id string) (*entities.User, error)

	// List returns a list of users with pagination support
	// offset: starting index, limit: maximum number of users to return
	// Returns the list of users and the total number of users
	List(ctx context.Context, offset, limit int) ([]*entities.User, int, error)

	// Update modifies an existing user's information
	Update(ctx context.Context, user *entities.User) (*entities.User, error)

	// UpdateStatus changes the status of a user (e.g., active/inactive)
	UpdateStatus(ctx context.Context, id string, status string) (*entities.User, error)

	// Delete removes a user from the repository by ID
	Delete(ctx context.Context, id string) error

	// GetByEmail retrieves a user by their unique email
	GetByEmail(ctx context.Context, email string) (*entities.User, error)

	// ExistsByEmail checks if a user with the given email already exists
	// Returns true if the email exists, false otherwise
	ExistsByEmail(ctx context.Context, email string) (bool, error)
}

// RefreshTokenRepository defines the interface for refresh token persistence
type RefreshTokenRepository interface {
	// SaveRefreshToken saves a refresh token for a user with a specific expiration.
	SaveRefreshToken(ctx context.Context, userID, token string, exp time.Time) error

	// GetRefreshToken retrieves a refresh token for a user
	GetRefreshToken(ctx context.Context, userID string) (string, error)

	// DeleteRefreshToken removes a refresh token (logout)
	DeleteRefreshToken(ctx context.Context, userID string) error
}
