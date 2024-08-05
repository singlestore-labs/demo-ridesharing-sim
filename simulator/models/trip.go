package models

import "time"

// Trip represents a single trip in the ridesharing simulation
type Trip struct {
	ID          string
	DriverID    string
	RiderID     string
	Status      string
	RequestTime time.Time
	AcceptTime  time.Time
	PickupTime  time.Time
	DropoffTime time.Time
	Fare        int
	Distance    float64
	PickupLat   float64
	PickupLong  float64
	DropoffLat  float64
	DropoffLong float64
	City        string
}
