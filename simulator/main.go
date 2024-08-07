package main

import (
	"simulator/api"
	"simulator/config"
	"simulator/database"
	"simulator/exporter"
	"simulator/service"
	"time"
)

func main() {
	service.LoadGeoData()
	database.Initialize()

	// riders, err := exporter.ImportRidersFromCSV("data/riders.csv")
	// if err != nil {
	// 	riders = service.GenerateRiders(config.NumRiders, "San Francisco")
	// }
	// if len(riders) != config.NumRiders {
	// 	riders = service.GenerateRiders(config.NumRiders, "San Francisco")
	// }

	// drivers, err := exporter.ImportDriversFromCSV("data/drivers.csv")
	// if err != nil {
	// 	drivers = service.GenerateDrivers(config.NumDrivers, "San Francisco")
	// }
	// if len(drivers) != config.NumDrivers {
	// 	drivers = service.GenerateDrivers(config.NumDrivers, "San Francisco")
	// }

	// exporter.ExportRidersToCSV(riders)
	// exporter.ExportDriversToCSV(drivers)

	riders := service.GenerateRiders(config.NumRiders, "San Francisco")
	drivers := service.GenerateDrivers(config.NumDrivers, "San Francisco")

	go func() {
		for _, rider := range riders {
			database.Local.Riders.Set(rider.ID, rider)
			go service.StartRiderLoop(rider.ID, "San Francisco")
			time.Sleep(time.Duration(config.Faker.IntBetween(1, 100)) * time.Millisecond)
		}
	}()

	go func() {
		for _, driver := range drivers {
			database.Local.Drivers.Set(driver.ID, driver)
			go service.StartDriverLoop(driver.ID, "San Francisco")
			time.Sleep(time.Duration(config.Faker.IntBetween(1, 100)) * time.Millisecond)
		}
	}()

	go func() {
		for {
			time.Sleep(1 * time.Minute)
			trips := service.GetTripsByStatus("completed")
			exporter.ExportTripsToCSV(trips)
		}
	}()

	api.StartServer()
}
