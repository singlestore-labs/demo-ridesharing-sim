package service

import (
	"math"
	"simulator/config"
	"simulator/database"
	"simulator/models"
	"slices"
	"time"
)

func FindTripsForRider(userID string) []models.Trip {
	trips := make([]models.Trip, 0)
	for _, trip := range database.Local.Trips {
		if trip.RiderID == userID {
			trips = append(trips, trip)
		}
	}
	return trips
}

func FindTripsForDriver(userID string) []models.Trip {
	trips := make([]models.Trip, 0)
	for _, trip := range database.Local.Trips {
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
	database.Local.Trips[trip.ID] = trip
}

func RequestRide(userID string, city string) {
	riderLocation := GetLocationForRider(userID)
	destLat, destLong := GenerateCoordinateWithinDistanceInCity(city, riderLocation.Latitude, riderLocation.Longitude, 5000)
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
}

func GetClosestRequest(lat, long float64) models.Trip {
	closestDistance := math.MaxFloat64
	var closestTrip models.Trip
	for _, trip := range database.Local.Trips {
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
	trip := database.Local.Trips[tripID]
	trip.DriverID = driverID
	trip.Status = "accepted"
	trip.AcceptTime = time.Now()
	UpsertTrip(trip)
}
