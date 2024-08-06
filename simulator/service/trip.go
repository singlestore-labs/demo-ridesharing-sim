package service

import (
	"simulator/config"
	"simulator/models"
)

func RequestRide(userID string, city string) models.Trip {
	trip := models.Trip{
		ID:      config.Faker.UUID().V4(),
		RiderID: userID,
	}
}
