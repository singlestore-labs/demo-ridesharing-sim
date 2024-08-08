package models

import "time"

// Trip represents a single trip in the ridesharing simulation
type Trip struct {
	ID       string `avro:"id"`
	DriverID string `avro:"driver_id"`
	RiderID  string `avro:"rider_id"`
	// Status can be "requested", "accepted", "en_route", "completed"
	Status      string    `avro:"status"`
	RequestTime time.Time `avro:"request_time"`
	AcceptTime  time.Time `avro:"accept_time"`
	PickupTime  time.Time `avro:"pickup_time"`
	DropoffTime time.Time `avro:"dropoff_time"`
	Fare        int       `avro:"fare"`
	Distance    float64   `avro:"distance"`
	PickupLat   float64   `avro:"pickup_lat"`
	PickupLong  float64   `avro:"pickup_long"`
	DropoffLat  float64   `avro:"dropoff_lat"`
	DropoffLong float64   `avro:"dropoff_long"`
	City        string    `avro:"city"`
}
