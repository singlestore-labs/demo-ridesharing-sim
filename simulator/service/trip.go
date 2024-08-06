package service

import (
	"math"
	"simulator/config"
	"simulator/database"
	"simulator/models"
	"slices"
	"time"
)

func GetAllTrips() []models.Trip {
	trips := make([]models.Trip, 0)
	for _, trip := range database.Local.Trips.Items() {
		trips = append(trips, trip)
	}
	return trips
}

func GetTrip(tripID string) models.Trip {
	trip, _ := database.Local.Trips.Get(tripID)
	return trip
}

func FindTripsForRider(userID string) []models.Trip {
	trips := make([]models.Trip, 0)
	for _, trip := range GetAllTrips() {
		if trip.RiderID == userID {
			trips = append(trips, trip)
		}
	}
	return trips
}

func FindTripsForDriver(userID string) []models.Trip {
	trips := make([]models.Trip, 0)
	for _, trip := range GetAllTrips() {
		if trip.DriverID == userID {
			trips = append(trips, trip)
		}
	}
	return trips
}

func GetLastTrip(trips []models.Trip) models.Trip {
	if len(trips) == 0 {
		return models.Trip{}
	}
	slices.SortFunc(trips, func(a, b models.Trip) int {
		return a.RequestTime.Compare(b.RequestTime)
	})
	return trips[len(trips)-1]
}

func UpsertTrip(trip models.Trip) {
	database.Local.Trips.Set(trip.ID, trip)
}

func RequestRide(userID string, city string) string {
	riderLocation := GetLocationForRider(userID)
	destLat, destLong := GenerateCoordinateWithinDistanceInCity(city, riderLocation.Latitude, riderLocation.Longitude, 500)
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
	for _, trip := range GetAllTrips() {
		if trip.Status == "requested" {
			distance := GetDistanceBetweenCoordinates(lat, long, trip.PickupLat, trip.PickupLong)
			if distance < closestDistance {
				closestDistance = distance
				closestTrip = trip
			}
		}
	}
	return closestTrip
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
	time.Sleep(1 * time.Second)
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
	time.Sleep(1 * time.Second)
	trip.Status = "completed"
	trip.DropoffTime = time.Now()
	UpsertTrip(trip)
}
