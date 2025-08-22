package postgresql

import (
	"database/sql"
	"errors"
	"go-chi-boilerplate/internal/core/entities"
	"go-chi-boilerplate/internal/core/repositories"
	"time"

	_ "github.com/lib/pq"
)

type userRepository struct {
	db *sql.DB
}

// Constructor
func NewUserRepository(db *sql.DB) repositories.UserRepository {
	return &userRepository{db: db}
}

// Create a new user
func (r *userRepository) Create(user *entities.User) (*entities.User, error) {
	query := `
	INSERT INTO users (first_name, last_name, email, role, password, status, created_at, updated_at)
	VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
	RETURNING id, created_at, updated_at
	`
	now := time.Now()
	err := r.db.QueryRow(query,
		user.FirstName, user.LastName, user.Email, user.Role,
		user.Password, user.Status, now, now,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// Get a user by ID
func (r *userRepository) GetByID(id string) (*entities.User, error) {
	user := &entities.User{}
	query := `
	SELECT id, first_name, last_name, email, role, password, status, created_at, updated_at
	FROM users WHERE id=$1
	`
	err := r.db.QueryRow(query, id).Scan(
		&user.ID, &user.FirstName, &user.LastName, &user.Email,
		&user.Role, &user.Password, &user.Status,
		&user.CreatedAt, &user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	}
	return user, err
}

// List users with pagination
func (r *userRepository) List(offset, limit int) ([]*entities.User, int, error) {
	query := `SELECT id, first_name, last_name, email, role, status, created_at, updated_at
			  FROM users ORDER BY created_at DESC LIMIT $1 OFFSET $2`
	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	users := []*entities.User{}
	for rows.Next() {
		u := &entities.User{}
		err := rows.Scan(&u.ID, &u.FirstName, &u.LastName, &u.Email, &u.Role, &u.Status, &u.CreatedAt, &u.UpdatedAt)
		if err != nil {
			return nil, 0, err
		}
		users = append(users, u)
	}

	// Count total users
	var total int
	err = r.db.QueryRow(`SELECT COUNT(*) FROM users`).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// Update a user
func (r *userRepository) Update(user *entities.User) (*entities.User, error) {
	query := `
	UPDATE users SET first_name=$1, last_name=$2, email=$3, role=$4, updated_at=$5
	WHERE id=$6 RETURNING created_at, updated_at
	`
	now := time.Now()
	err := r.db.QueryRow(query, user.FirstName, user.LastName, user.Email, user.Role, now, user.ID).
		Scan(&user.CreatedAt, &user.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	}
	return user, err
}

// Update user status
func (r *userRepository) UpdateStatus(id string, status string) (*entities.User, error) {
	query := `
	UPDATE users SET status=$1, updated_at=$2 WHERE id=$3
	RETURNING first_name, last_name, email, role, created_at, updated_at
	`
	user := &entities.User{}
	err := r.db.QueryRow(query, status, time.Now(), id).Scan(
		&user.FirstName, &user.LastName, &user.Email, &user.Role, &user.CreatedAt, &user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	}
	user.ID = id
	user.Status = status
	return user, err
}

// Delete a user
func (r *userRepository) Delete(id string) error {
	query := `DELETE FROM users WHERE id=$1`
	res, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}
	count, _ := res.RowsAffected()
	if count == 0 {
		return errors.New("user not found")
	}
	return nil
}

// Get user with email id
func (r *userRepository) GetByEmail(email string) (*entities.User, error) {
	var user entities.User
	query := "SELECT id, first_name, last_name, email, role, password, status FROM users WHERE email=$1"
	err := r.db.QueryRow(query, email).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Role,
		&user.Password,
		&user.Status,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// ExistsByEmail checks if a user with the given email already exists
func (r *userRepository) ExistsByEmail(email string) (bool, error) {
	var count int
	err := r.db.QueryRow("SELECT COUNT(1) FROM users WHERE email = $1", email).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
