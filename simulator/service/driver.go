package service

import (
	"simulator/config"
	"simulator/models"
	"time"
)

func GenerateDriver(city string) models.Driver {
	lat, long := GenerateCoordinateInCity(city)
	driver := models.Driver{
		ID:          config.Faker.UUID().V4(),
		FirstName:   config.Faker.Person().FirstName(),
		LastName:    config.Faker.Person().LastName(),
		Email:       config.Faker.Internet().Email(),
		PhoneNumber: config.Faker.Phone().Number(),
		DateOfBirth: config.Faker.Time().TimeBetween(time.Now().AddDate(-30, 0, 0), time.Now()),
		CreatedAt:   time.Now(),
	}
	driver.Location = models.Location{
		UserID:    driver.ID,
		Latitude:  lat,
		Longitude: long,
	}
	return driver
}

func GenerateDrivers(numDrivers int, city string) []models.Driver {
	drivers := make([]models.Driver, numDrivers)
	for i := 0; i < numDrivers; i++ {
		drivers[i] = GenerateDriver(city)
	}
	return drivers
}
