package main

import (
	"fmt"
	"net/http"
	"simulator/config"
	"simulator/exporter"
	"simulator/service"
	"strings"
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

	lat, long := service.GenerateCoordinateInCity("San Francisco")
	fmt.Println(lat, long)

	cordinates := make([]string, 0)
	for i := 0; i < 1000; i++ {
		// lat, long := service.GenerateCoordinateWithinDistanceInCity("San Francisco", lat, long, 1000)
		lat, long = service.GenerateCoordinateInCity("San Francisco")
		cordinates = append(cordinates, fmt.Sprintf("{\"latitude\": %f, \"longitude\": %f}", lat, long))
	}

	http.HandleFunc("/coordinates", func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// Handle preflight OPTIONS request
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		response := fmt.Sprintf("[%s]", strings.Join(cordinates, ", "))
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(response))
	})

	fmt.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Error starting server: %s\n", err)
	}
}
