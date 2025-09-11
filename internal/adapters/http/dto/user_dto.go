package dto

import "time"

// CreateUserRequest represents the payload to create a new user
type CreateUserRequest struct {
	FirstName string `json:"first_name" example:"John" validate:"required" description:"User's first name"`
	LastName  string `json:"last_name" example:"Doe" validate:"required" description:"User's last name"`
	Email     string `json:"email" example:"john.doe@example.com" validate:"required,email" description:"User's email address"`
	Role      string `json:"role" example:"admin" validate:"required" description:"Role assigned to the user"`
	Password  string `json:"password" example:"P@ssw0rd" validate:"required,min=8" description:"Password for the user account"`
}

// UpdateUserRequest represents the payload to update an existing user
type UpdateUserRequest struct {
	FirstName *string `json:"first_name,omitempty" description:"Updated first name of the user"`
	LastName  *string `json:"last_name,omitempty" description:"Updated last name of the user"`
	Email     *string `json:"email,omitempty" description:"Updated email address of the user"`
	Role      *string `json:"role,omitempty" description:"Updated role of the user"`
}

// ChangeUserStatusRequest represents the payload to change a user's status
type ChangeUserStatusRequest struct {
	Status string `json:"status" example:"active" validate:"required,oneof=active inactive" description:"New status of the user"`
}

// UserResponse represents the user object returned in API responses
type UserResponse struct {
	ID        string    `json:"id" description:"Unique identifier of the user"`
	FirstName string    `json:"first_name" description:"First name of the user"`
	LastName  string    `json:"last_name" description:"Last name of the user"`
	Email     string    `json:"email" description:"Email address of the user"`
	Role      string    `json:"role" description:"Role assigned to the user"`
	Status    string    `json:"status" description:"Current status of the user"`
	CreatedAt time.Time `json:"created_at" description:"Timestamp when the user was created"`
	UpdatedAt time.Time `json:"updated_at" description:"Timestamp when the user was last updated"`
}

// ListUsersResponse represents the response for listing users
type ListUsersResponse struct {
	Users []*UserResponse `json:"users" description:"List of users"`
	Total int             `json:"total" description:"Total number of users"`
}
