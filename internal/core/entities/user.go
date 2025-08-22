package entities

import "time"

type User struct {
	ID        string
	FirstName string
	LastName  string
	Email     string
	Role      string
	Password  string
	Status    string
	CreatedAt time.Time
	UpdatedAt time.Time
}
