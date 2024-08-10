package service

import (
	"fmt"
	"simulator/config"
	"simulator/database"
	"simulator/models"
	"time"
)

// ================================
//  SIMULATION FUNCTIONS
// ================================

func StartRiderLoop(userID string, city string) {
	for {
		UpdateStatusForRider(userID, "idle")
		initLat, initLong := GenerateCoordinateInCity(city)
		UpdateLocationForRider(userID, models.Location{
			Latitude:  initLat,
			Longitude: initLong,
		})
		sleepTime := time.Duration(config.Faker.IntBetween(500, 20000)) * time.Millisecond
		fmt.Printf("Rider %s is idle for %s\n", userID, sleepTime)
		// time.Sleep(sleepTime)
		fmt.Println("Requesting ride for rider...", userID)
		tripID := RequestRide(userID, city)
		if tripID == "" {
			fmt.Printf("Rider %s failed to request ride\n", userID)
			continue
		}
		UpdateStatusForRider(userID, "requested")
		fmt.Printf("Rider %s requested ride %s\n", userID, tripID)
		for GetTrip(tripID).Status != "accepted" {
			time.Sleep(100 * time.Millisecond)
		}
		UpdateStatusForRider(userID, "waiting")
		for GetTrip(tripID).Status != "en_route" {
			time.Sleep(100 * time.Millisecond)
		}
		UpdateStatusForRider(userID, "in_progress")
		for GetTrip(tripID).Status != "completed" {
			time.Sleep(100 * time.Millisecond)
		}
		fmt.Printf("Rider %s completed trip %s\n\n", userID, tripID)
	}
}

func GenerateRider(city string) models.Rider {
	lat, long := GenerateCoordinateInCity(city)
	rider := models.Rider{
		ID:          config.Faker.UUID().V4(),
		FirstName:   config.Faker.Person().FirstName(),
		LastName:    config.Faker.Person().LastName(),
		Email:       config.Faker.Internet().Email(),
		PhoneNumber: config.Faker.Phone().Number(),
		DateOfBirth: config.Faker.Time().TimeBetween(time.Now().AddDate(-30, 0, 0), time.Now()),
		CreatedAt:   time.Now(),
	}
	rider.Location = models.Location{
		UserID:    rider.ID,
		Latitude:  lat,
		Longitude: long,
		City:      city,
		Timestamp: time.Now(),
	}
	return rider
}

func GenerateRiders(numRiders int, city string) []models.Rider {
	riders := make([]models.Rider, numRiders)
	for i := 0; i < numRiders; i++ {
		riders[i] = GenerateRider(city)
	}
	return riders
}

// ================================
//  LOCAL DATABASE FUNCTIONS
// ================================

func GetAllRiders() []models.Rider {
	riders := make([]models.Rider, 0)
	for _, rider := range database.Local.Riders.Items() {
		riders = append(riders, rider)
	}
	return riders
}

func GetRider(userID string) models.Rider {
	rider, ok := database.Local.Riders.Get(userID)
	if !ok {
		return models.Rider{}
	}
	return rider
}

func GetLocationForRider(userID string) models.Location {
	rider := GetRider(userID)
	return rider.Location
}

func UpdateLocationForRider(userID string, location models.Location) {
	rider := GetRider(userID)
	rider.Location = location
	rider.Location.UserID = userID
	rider.Location.Timestamp = time.Now()
	database.Local.Riders.Set(userID, rider)
}

func UpdateStatusForRider(userID string, status string) {
	rider := GetRider(userID)
	if rider.ID != "" {
		rider.Status = status
		database.Local.Riders.Set(userID, rider)
	}
}
