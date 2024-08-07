package models

import "time"

type Location struct {
	UserID    string    `json:"id"`
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
	City      string    `json:"city"`
	Timestamp time.Time `json:"timestamp"`
}
