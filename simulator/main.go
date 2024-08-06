package main

import (
	"simulator/api"
	"simulator/config"
	"simulator/database"
	"simulator/service"
	"time"
)

func main() {
	service.LoadGeoData()
	database.Initialize()

	api.NewResourceArrayEndpoint("/riders", api.SerializeMapToJSONArray(database.Local.Riders))
	api.NewResourceArrayEndpoint("/drivers", api.SerializeMapToJSONArray(database.Local.Drivers))
	api.NewResourceArrayEndpoint("/trips", api.SerializeMapToJSONArray(database.Local.Trips))
	api.Serve()

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

	for _, rider := range riders {
		database.Local.Riders[rider.ID] = rider
	}
	for _, driver := range drivers {
		database.Local.Drivers[driver.ID] = driver
	}

	service.RequestRide(riders[0].ID, "San Francisco")
	time.Sleep(5 * time.Second)
	service.AcceptRide(service.GetLastTrip(service.FindTripsForRider(riders[0].ID)).ID, drivers[0].ID)
}
