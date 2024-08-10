package service

import (
	"math"
	"simulator/config"
	"simulator/database"
	"simulator/exporter"
	"simulator/model"
	"time"
)

// ================================
//  SIMULATION FUNCTIONS
// ================================

func RequestRide(userID string, city string) string {
	lat, long := GetLocationForRider(userID)
	tripDistance := config.Faker.RandomFloat(1, 100, 16000)
	destLat, destLong := GenerateCoordinateWithinDistanceInCity(city, lat, long, tripDistance)
	if destLat == 0 && destLong == 0 {
		return ""
	}
	trip := model.Trip{
		ID:          config.Faker.UUID().V4(),
		RiderID:     userID,
		Status:      "requested",
		RequestTime: time.Now(),
		City:        city,
		PickupLat:   lat,
		PickupLong:  long,
		DropoffLat:  destLat,
		DropoffLong: destLong,
		Distance:    GetDistanceBetweenCoordinates(lat, long, destLat, destLong),
	}
	UpsertTrip(trip)
	return trip.ID
}

func GetClosestRequest(lat, long float64) model.Trip {
	closestDistance := math.MaxFloat64
	var closestTrip model.Trip
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
	lat, long := GetLocationForDriver(trip.DriverID)
	path := GenerateMiddleCoordinates(lat, long, trip.PickupLat, trip.PickupLong, 10)
	for _, point := range path {
		time.Sleep(100 * time.Millisecond)
		UpdateLocationForDriver(trip.DriverID, point[0], point[1])
	}
	UpdateLocationForDriver(trip.DriverID, trip.PickupLat, trip.PickupLong)
	// pickup rider
	time.Sleep(time.Duration(config.Faker.IntBetween(200, 3000)) * time.Millisecond)
	trip.Status = "en_route"
	trip.PickupTime = time.Now()
	UpsertTrip(trip)
	// driver to dropoff
	path = GenerateMiddleCoordinates(trip.PickupLat, trip.PickupLong, trip.DropoffLat, trip.DropoffLong, 10)
	for _, point := range path {
		time.Sleep(100 * time.Millisecond)
		UpdateLocationForDriver(trip.DriverID, point[0], point[1])
		UpdateLocationForRider(trip.RiderID, point[0], point[1])
	}
	UpdateLocationForDriver(trip.DriverID, trip.DropoffLat, trip.DropoffLong)
	UpdateLocationForRider(trip.RiderID, trip.DropoffLat, trip.DropoffLong)
	// dropoff rider
	time.Sleep(time.Duration(config.Faker.IntBetween(200, 3000)) * time.Millisecond)
	trip.Status = "completed"
	trip.DropoffTime = time.Now()
	UpsertTrip(trip)
}

// ================================
//  LOCAL DATABASE FUNCTIONS
// ================================

func GetAllTrips() []model.Trip {
	trips := make([]model.Trip, 0)
	for _, trip := range database.Local.Trips.Items() {
		trips = append(trips, trip)
	}
	return trips
}

func GetTripsByStatus(status string) []model.Trip {
	trips := make([]model.Trip, 0)
	for _, trip := range database.Local.Trips.Items() {
		if trip.Status == status {
			trips = append(trips, trip)
		}
	}
	return trips
}

func GetTrip(tripID string) model.Trip {
	trip, ok := database.Local.Trips.Get(tripID)
	if !ok {
		return model.Trip{}
	}
	return trip
}

func UpsertTrip(trip model.Trip) {
	database.Local.Trips.Set(trip.ID, trip)
	go exporter.KafkaProduceTrip(trip)
}
