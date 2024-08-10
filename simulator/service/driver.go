package service

import (
	"log"
	"simulator/config"
	"simulator/database"
	"simulator/exporter"
	"simulator/model"
	"time"
)

// ================================
//  SIMULATION FUNCTIONS
// ================================

func StartDriverLoop(userID string, city string) {
	initLat, initLong := GenerateCoordinateInCity(city)
	UpdateLocationForDriver(userID, model.Location{
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
				// log.Printf("Driver %s waiting for request\n", userID)
				time.Sleep(100 * time.Millisecond)
				request = GetClosestRequest(driverLocation.Latitude, driverLocation.Longitude)
			}
			accepted = TryAcceptRide(request.ID, userID)
			if !accepted {
				request.ID = ""
			}
		}
		UpdateStatusForDriver(userID, "in_progress")
		log.Printf("Driver %s accepted request %s\n", userID, request.ID)
		StartTripLoop(request.ID)
		log.Printf("Driver %s completed trip %s\n\n", userID, request.ID)
	}
}

func GenerateDriver(city string) model.Driver {
	lat, long := GenerateCoordinateInCity(city)
	driver := model.Driver{
		ID:          config.Faker.UUID().V4(),
		FirstName:   config.Faker.Person().FirstName(),
		LastName:    config.Faker.Person().LastName(),
		Email:       config.Faker.Internet().Email(),
		PhoneNumber: config.Faker.Phone().Number(),
		DateOfBirth: config.Faker.Time().TimeBetween(time.Now().AddDate(-30, 0, 0), time.Now()),
		CreatedAt:   time.Now(),
	}
	driver.Location = model.Location{
		UserID:    driver.ID,
		Latitude:  lat,
		Longitude: long,
		City:      city,
		Timestamp: time.Now(),
	}
	return driver
}

func GenerateDrivers(numDrivers int, city string) []model.Driver {
	drivers := make([]model.Driver, numDrivers)
	for i := 0; i < numDrivers; i++ {
		drivers[i] = GenerateDriver(city)
	}
	return drivers
}

// ================================
//  LOCAL DATABASE FUNCTIONS
// ================================

func GetAllDrivers() []model.Driver {
	drivers := make([]model.Driver, 0)
	for _, driver := range database.Local.Drivers.Items() {
		drivers = append(drivers, driver)
	}
	return drivers
}

func GetDriver(userID string) model.Driver {
	driver, ok := database.Local.Drivers.Get(userID)
	if !ok {
		return model.Driver{}
	}
	return driver
}

func GetLocationForDriver(userID string) model.Location {
	driver := GetDriver(userID)
	return driver.Location
}

func UpdateLocationForDriver(userID string, location model.Location) {
	driver := GetDriver(userID)
	if driver.ID == "" {
		return
	}
	driver.Location.UserID = userID
	driver.Location.Latitude = location.Latitude
	driver.Location.Longitude = location.Longitude
	driver.Location.Timestamp = time.Now()
	database.Local.Drivers.Set(userID, driver)
	exporter.KafkaProduceDriver(driver)
}

func UpdateStatusForDriver(userID string, status string) {
	driver := GetDriver(userID)
	if driver.ID == "" {
		return
	}
	driver.Status = status
	database.Local.Drivers.Set(userID, driver)
	exporter.KafkaProduceDriver(driver)
}
