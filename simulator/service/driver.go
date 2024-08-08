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

func StartDriverLoop(userID string, city string) {
	initLat, initLong := GenerateCoordinateInCity(city)
	UpdateLocationForDriver(userID, models.Location{
		Latitude:  initLat,
		Longitude: initLong,
	})
	for {
		UpdateStatusForDriver(userID, "available")
		sleepTime := time.Duration(config.Faker.IntBetween(500, 2000)) * time.Millisecond
		time.Sleep(sleepTime)
		driverLocation := GetLocationForDriver(userID)
		request := GetClosestRequest(driverLocation.Latitude, driverLocation.Longitude)
		accepted := false
		for !accepted {
			for request.ID == "" {
				// fmt.Printf("Driver %s waiting for request\n", userID)
				time.Sleep(100 * time.Millisecond)
				request = GetClosestRequest(driverLocation.Latitude, driverLocation.Longitude)
			}
			accepted = TryAcceptRide(request.ID, userID)
			if !accepted {
				request.ID = ""
			}
		}
		UpdateStatusForDriver(userID, "in_progress")
		fmt.Printf("Driver %s accepted request %s\n", userID, request.ID)
		StartTripLoop(request.ID)
		fmt.Printf("Driver %s completed trip %s\n\n", userID, request.ID)
	}
}

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
		City:      city,
		Timestamp: time.Now(),
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

// ================================
//  LOCAL DATABASE FUNCTIONS
// ================================

func GetAllDrivers() []models.Driver {
	drivers := make([]models.Driver, 0)
	for _, driver := range database.Local.Drivers.Items() {
		drivers = append(drivers, driver)
	}
	return drivers
}

func GetDriversByStatus(status string) []models.Driver {
	drivers := make([]models.Driver, 0)
	for _, driver := range database.Local.Drivers.Items() {
		if driver.Status == status {
			drivers = append(drivers, driver)
		}
	}
	return drivers
}

func GetDriver(userID string) models.Driver {
	driver, ok := database.Local.Drivers.Get(userID)
	if !ok {
		return models.Driver{}
	}
	return driver
}

func GetLocationForDriver(userID string) models.Location {
	driver := GetDriver(userID)
	return driver.Location
}

func UpdateLocationForDriver(userID string, location models.Location) {
	driver := GetDriver(userID)
	driver.Location = location
	driver.Location.UserID = userID
	driver.Location.Timestamp = time.Now()
	database.Local.Drivers.Set(userID, driver)
}

func UpdateStatusForDriver(userID string, status string) {
	driver := GetDriver(userID)
	if driver.ID != "" {
		driver.Status = status
		database.Local.Drivers.Set(userID, driver)
	}
}
