package main

import (
	"simulator/api"
	"simulator/config"
	"simulator/exporter"
	"simulator/service"
)

func main() {
	service.LoadGeoData()

	riders, err := exporter.ImportRidersFromCSV("data/riders.csv")
	if err != nil {
		riders = service.GenerateRiders(config.NumRiders)
	}
	if len(riders) != config.NumRiders {
		riders = service.GenerateRiders(config.NumRiders)
	}

	drivers, err := exporter.ImportDriversFromCSV("data/drivers.csv")
	if err != nil {
		drivers = service.GenerateDrivers(config.NumDrivers)
	}
	if len(drivers) != config.NumDrivers {
		drivers = service.GenerateDrivers(config.NumDrivers)
	}

	exporter.ExportRidersToCSV(riders)
	exporter.ExportDriversToCSV(drivers)

	cordinates := make([]string, 0)

	api.NewResourceEndpoint("/coordinates", cordinates)
	api.Serve()
}
