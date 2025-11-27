package entity

import "time"

// User represents a user entity belonging to a company.
type User struct {
	ID           int64
	CompanyID    int64
	Name         string
	Email        string
	PasswordHash string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
