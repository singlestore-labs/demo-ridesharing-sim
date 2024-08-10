package main

import (
	"simulator/config"
	"simulator/database"
	"simulator/exporter"
	"simulator/service"
	"time"
)

func main() {
	config.Verify()
	service.LoadGeoData()
	database.InitializeLocal()
	exporter.InitializeKafkaClient()

	riders := service.GenerateRiders(config.NumRiders, "San Francisco")
	drivers := service.GenerateDrivers(config.NumDrivers, "San Francisco")

	go func() {
		for _, rider := range riders {
			database.Local.Riders.Set(rider.ID, rider)
			go service.StartRiderLoop(rider.ID, "San Francisco")
			time.Sleep(time.Duration(config.Faker.IntBetween(1, 1000)) * time.Millisecond)
		}
	}()

	go func() {
		for _, driver := range drivers {
			database.Local.Drivers.Set(driver.ID, driver)
			go service.StartDriverLoop(driver.ID, "San Francisco")
			time.Sleep(time.Duration(config.Faker.IntBetween(1, 1000)) * time.Millisecond)
		}
	}()

	// Keep the main thread alive
	select {}
}
