package models

import (
	"time"
)

// Rider represents a rider account in a ridesharing app
type Rider struct {
	ID          string    `json:"id"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	Email       string    `json:"email"`
	PhoneNumber string    `json:"phone_number"`
	DateOfBirth time.Time `json:"date_of_birth"`
	CreatedAt   time.Time `json:"created_at"`
	Location    Location  `json:"location"`
	Status      string    `json:"status"`
}
