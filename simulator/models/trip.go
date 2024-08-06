package models

import "time"

// Trip represents a single trip in the ridesharing simulation
type Trip struct {
	ID       string `json:"id"`
	DriverID string `json:"driver_id"`
	RiderID  string `json:"rider_id"`
	// Status can be "requested", "accepted", "en_route", "completed"
	Status      string    `json:"status"`
	RequestTime time.Time `json:"request_time"`
	AcceptTime  time.Time `json:"accept_time"`
	PickupTime  time.Time `json:"pickup_time"`
	DropoffTime time.Time `json:"dropoff_time"`
	Fare        int       `json:"fare"`
	Distance    float64   `json:"distance"`
	PickupLat   float64   `json:"pickup_lat"`
	PickupLong  float64   `json:"pickup_long"`
	DropoffLat  float64   `json:"dropoff_lat"`
	DropoffLong float64   `json:"dropoff_long"`
	City        string    `json:"city"`
}
