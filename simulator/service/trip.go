package service

import (
	"fmt"
	"math"
	"simulator/config"
	"simulator/database"
	"simulator/models"
	"time"
)

// ================================
//  SIMULATION FUNCTIONS
// ================================

func RequestRide(userID string, city string) string {
	riderLocation := GetLocationForRider(userID)
	tripDistance := config.Faker.RandomFloat(1, 100, 16000)
	destLat, destLong := GenerateCoordinateWithinDistanceInCity(city, riderLocation.Latitude, riderLocation.Longitude, tripDistance)
	trip := models.Trip{
		ID:          config.Faker.UUID().V4(),
		RiderID:     userID,
		Status:      "requested",
		RequestTime: time.Now(),
		City:        city,
		PickupLat:   riderLocation.Latitude,
		PickupLong:  riderLocation.Longitude,
		DropoffLat:  destLat,
		DropoffLong: destLong,
		Distance:    GetDistanceBetweenCoordinates(riderLocation.Latitude, riderLocation.Longitude, destLat, destLong),
	}
	UpsertTrip(trip)
	return trip.ID
}

func GetClosestRequest(lat, long float64) models.Trip {
	closestDistance := math.MaxFloat64
	var closestTrip models.Trip
	for _, trip := range GetTripsByStatus("requested") {
		distance := GetDistanceBetweenCoordinates(lat, long, trip.PickupLat, trip.PickupLong)
		if distance < closestDistance {
			closestDistance = distance
			closestTrip = trip
		}
	}
	return closestTrip
}

func TryAcceptRide(tripID string, driverID string) bool {
	database.Local.AcceptMutex.Lock()
	defer database.Local.AcceptMutex.Unlock()
	if GetTrip(tripID).Status == "requested" {
		AcceptRide(tripID, driverID)
		return true
	}
	return false
}

func AcceptRide(tripID string, driverID string) {
	trip := GetTrip(tripID)
	trip.DriverID = driverID
	trip.Status = "accepted"
	trip.AcceptTime = time.Now()
	UpsertTrip(trip)
}

func StartTripLoop(tripID string) {
	trip := GetTrip(tripID)
	// driver to pickup
	driverLocation := GetLocationForDriver(trip.DriverID)
	path := GenerateMiddleCoordinates(driverLocation.Latitude, driverLocation.Longitude, trip.PickupLat, trip.PickupLong, 10)
	for _, point := range path {
		time.Sleep(100 * time.Millisecond)
		UpdateLocationForDriver(trip.DriverID, models.Location{
			Latitude:  point[0],
			Longitude: point[1],
		})
	}
	UpdateLocationForDriver(trip.DriverID, models.Location{
		Latitude:  trip.PickupLat,
		Longitude: trip.PickupLong,
	})
	// pickup rider
	time.Sleep(time.Duration(config.Faker.IntBetween(200, 3000)) * time.Millisecond)
	trip.Status = "en_route"
	trip.PickupTime = time.Now()
	UpsertTrip(trip)
	// driver to dropoff
	path = GenerateMiddleCoordinates(trip.PickupLat, trip.PickupLong, trip.DropoffLat, trip.DropoffLong, 10)
	for _, point := range path {
		time.Sleep(100 * time.Millisecond)
		UpdateLocationForDriver(trip.DriverID, models.Location{
			Latitude:  point[0],
			Longitude: point[1],
		})
		UpdateLocationForRider(trip.RiderID, models.Location{
			Latitude:  point[0],
			Longitude: point[1],
		})
	}
	UpdateLocationForDriver(trip.DriverID, models.Location{
		Latitude:  trip.DropoffLat,
		Longitude: trip.DropoffLong,
	})
	UpdateLocationForRider(trip.RiderID, models.Location{
		Latitude:  trip.DropoffLat,
		Longitude: trip.DropoffLong,
	})
	// dropoff rider
	time.Sleep(time.Duration(config.Faker.IntBetween(200, 3000)) * time.Millisecond)
	trip.Status = "completed"
	trip.DropoffTime = time.Now()
	UpsertTrip(trip)
}

// ================================
//  LOCAL DATABASE FUNCTIONS
// ================================

func GetAllTrips() []models.Trip {
	trips := make([]models.Trip, 0)
	for _, trip := range database.Local.Trips.Items() {
		trips = append(trips, trip)
	}
	return trips
}

func GetTripsByStatus(status string) []models.Trip {
	startTime := time.Now()
	defer func() {
		executionTime := time.Since(startTime)
		fmt.Printf("GetTripsByStatus execution time: %v\n", executionTime)
	}()

	trips := make([]models.Trip, 0)
	for _, trip := range database.Local.Trips.Items() {
		if trip.Status == status {
			trips = append(trips, trip)
		}
	}
	return trips
}

func GetTrip(tripID string) models.Trip {
	trip, ok := database.Local.Trips.Get(tripID)
	if !ok {
		return models.Trip{}
	}
	return trip
}

func UpsertTrip(trip models.Trip) {
	database.Local.Trips.Set(trip.ID, trip)
}
