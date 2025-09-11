package dto

// ProfileResponse represents the JSON response for the /profile endpoint
type ProfileResponse struct {
	UserID string `json:"user_id" example:"123e4567-e89b-12d3-a456-426614174000" description:"Unique user identifier from JWT"`

	Role string `json:"role" example:"admin" description:"Role of the user from JWT"`

	FirstName string `json:"first_name" example:"John" description:"User's first name"`

	LastName string `json:"last_name" example:"Doe" description:"User's last name"`

	Email string `json:"email" example:"john.doe@example.com" description:"User's registered email address"`
}
