package service

import (
	"simulator/config"
	"simulator/database"
	"simulator/models"
	"time"
)

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

func GetLocationForRider(userID string) models.Location {
	return database.Local.Riders[userID].Location
}

func UpdateLocationForRider(userID string, location models.Location) {
	rider := database.Local.Riders[userID]
	rider.Location = location
	rider.Location.UserID = userID
	rider.Location.Timestamp = time.Now()
	database.Local.Riders[userID] = rider
}
