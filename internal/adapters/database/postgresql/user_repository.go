package postgresql

import (
	"context"
	"database/sql"
	"errors"
	"go-chi-boilerplate/internal/core/entities"
	"go-chi-boilerplate/internal/core/repositories"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	_ "github.com/lib/pq"
)

// SQL Queries as constants
const (
	createUserQuery = `
	INSERT INTO users (first_name, last_name, email, role, password, status, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	RETURNING id, created_at, updated_at
	`

	getUserByIDQuery = `
	SELECT id, first_name, last_name, email, role, password, status, created_at, updated_at
	FROM users WHERE id = $1
	`

	listUsersQuery = `
	SELECT id, first_name, last_name, email, role, status, created_at, updated_at
	FROM users ORDER BY created_at DESC LIMIT $1 OFFSET $2
	`

	countUsersQuery = `SELECT COUNT(*) FROM users`

	updateUserQuery = `
	UPDATE users SET first_name = $1, last_name = $2, email = $3, role = $4, updated_at = $5
	WHERE id = $6 RETURNING created_at, updated_at
	`

	updateUserStatusQuery = `
	UPDATE users SET status = $1, updated_at = $2 WHERE id = $3
	RETURNING first_name, last_name, email, role, created_at, updated_at
	`

	deleteUserQuery = `DELETE FROM users WHERE id = $1`

	getUserByEmailQuery = `
	SELECT id, first_name, last_name, email, role, password, status, created_at, updated_at
	FROM users WHERE email = $1
	`

	existsByEmailQuery = `SELECT COUNT(1) FROM users WHERE email = $1`
)

type userRepository struct {
	db     *sql.DB
	tracer trace.Tracer
}

// NewUserRepository creates a new user repository with tracing
func NewUserRepository(db *sql.DB) repositories.UserRepository {
	return &userRepository{
		db:     db,
		tracer: otel.Tracer("postgresql.userRepository"),
	}
}

// Create a new user
func (r *userRepository) Create(ctx context.Context, user *entities.User) (*entities.User, error) {
	ctx, span := r.tracer.Start(ctx, "CreateUser")
	defer span.End()

	now := time.Now()
	err := r.db.QueryRowContext(ctx, createUserQuery,
		user.FirstName, user.LastName, user.Email, user.Role,
		user.Password, user.Status, now, now,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	span.SetAttributes(attribute.String("user.id", user.ID))
	return user, nil
}

// GetByID gets a user by ID
func (r *userRepository) GetByID(ctx context.Context, id string) (*entities.User, error) {
	ctx, span := r.tracer.Start(ctx, "GetUserByID")
	defer span.End()

	user := &entities.User{}
	err := r.db.QueryRowContext(ctx, getUserByIDQuery, id).Scan(
		&user.ID, &user.FirstName, &user.LastName, &user.Email,
		&user.Role, &user.Password, &user.Status,
		&user.CreatedAt, &user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		err = errors.New("user not found")
		span.RecordError(err)
		return nil, err
	}
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	span.SetAttributes(attribute.String("user.id", id))
	return user, nil
}

// List users with pagination
func (r *userRepository) List(ctx context.Context, offset, limit int) ([]*entities.User, int, error) {
	ctx, span := r.tracer.Start(ctx, "ListUsers")
	defer span.End()

	rows, err := r.db.QueryContext(ctx, listUsersQuery, limit, offset)
	if err != nil {
		span.RecordError(err)
		return nil, 0, err
	}
	defer rows.Close()

	users := []*entities.User{}
	for rows.Next() {
		u := &entities.User{}
		err := rows.Scan(&u.ID, &u.FirstName, &u.LastName, &u.Email, &u.Role, &u.Status, &u.CreatedAt, &u.UpdatedAt)
		if err != nil {
			span.RecordError(err)
			return nil, 0, err
		}
		users = append(users, u)
	}

	var total int
	err = r.db.QueryRowContext(ctx, countUsersQuery).Scan(&total)
	if err != nil {
		span.RecordError(err)
		return nil, 0, err
	}

	span.SetAttributes(attribute.Int("user.count", total))
	return users, total, nil
}

// Update a user
func (r *userRepository) Update(ctx context.Context, user *entities.User) (*entities.User, error) {
	ctx, span := r.tracer.Start(ctx, "UpdateUser")
	defer span.End()

	now := time.Now()
	err := r.db.QueryRowContext(ctx, updateUserQuery,
		user.FirstName, user.LastName, user.Email, user.Role, now, user.ID,
	).Scan(&user.CreatedAt, &user.UpdatedAt)
	if err == sql.ErrNoRows {
		err = errors.New("user not found")
		span.RecordError(err)
		return nil, err
	}
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	span.SetAttributes(attribute.String("user.id", user.ID))
	return user, nil
}

// Update user status
func (r *userRepository) UpdateStatus(ctx context.Context, id string, status string) (*entities.User, error) {
	ctx, span := r.tracer.Start(ctx, "UpdateUserStatus")
	defer span.End()

	user := &entities.User{}
	err := r.db.QueryRowContext(ctx, updateUserStatusQuery, status, time.Now(), id).Scan(
		&user.FirstName, &user.LastName, &user.Email, &user.Role, &user.CreatedAt, &user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		err = errors.New("user not found")
		span.RecordError(err)
		return nil, err
	}
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	user.ID = id
	user.Status = status
	span.SetAttributes(attribute.String("user.id", id), attribute.String("user.status", status))
	return user, nil
}

// Delete a user
func (r *userRepository) Delete(ctx context.Context, id string) error {
	ctx, span := r.tracer.Start(ctx, "DeleteUser")
	defer span.End()

	res, err := r.db.ExecContext(ctx, deleteUserQuery, id)
	if err != nil {
		span.RecordError(err)
		return err
	}

	count, _ := res.RowsAffected()
	if count == 0 {
		err = errors.New("user not found")
		span.RecordError(err)
		return err
	}

	span.SetAttributes(attribute.String("user.id", id))
	return nil
}

// GetByEmail gets a user with email id
func (r *userRepository) GetByEmail(ctx context.Context, email string) (*entities.User, error) {
	ctx, span := r.tracer.Start(ctx, "GetUserByEmail")
	defer span.End()

	user := &entities.User{}
	err := r.db.QueryRowContext(ctx, getUserByEmailQuery, email).Scan(
		&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Role,
		&user.Password, &user.Status, &user.CreatedAt, &user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		err = errors.New("user not found")
		span.RecordError(err)
		return nil, err
	}
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	span.SetAttributes(attribute.String("user.id", user.ID), attribute.String("user.email", user.Email))
	return user, nil
}

// ExistsByEmail checks if a user with the given email already exists
func (r *userRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	ctx, span := r.tracer.Start(ctx, "ExistsUserByEmail")
	defer span.End()

	var count int
	err := r.db.QueryRowContext(ctx, existsByEmailQuery, email).Scan(&count)
	if err != nil {
		span.RecordError(err)
		return false, err
	}

	span.SetAttributes(attribute.String("user.email", email), attribute.Bool("user.exists", count > 0))
	return count > 0, nil
}
