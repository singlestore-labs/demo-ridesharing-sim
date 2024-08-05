package exporter

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"simulator/models"
	"time"
)

func ExportRidersToCSV(riders []models.Rider) {
	// Create CSV file for riders
	riderFile, err := os.Create("data/riders.csv")
	if err != nil {
		log.Fatal("Could not create riders.csv file", err)
	}
	defer riderFile.Close()

	riderWriter := csv.NewWriter(riderFile)
	defer riderWriter.Flush()

	// Write header for riders CSV
	riderWriter.Write([]string{"id", "first_name", "last_name", "email", "phone_number", "date_of_birth", "created_at"})

	// Generate and write riders
	for _, rider := range riders {
		riderWriter.Write([]string{
			rider.ID,
			rider.FirstName,
			rider.LastName,
			rider.Email,
			rider.PhoneNumber,
			rider.DateOfBirth.Format(time.RFC3339),
			rider.CreatedAt.Format(time.RFC3339),
		})
	}
}

func ImportRidersFromCSV(filePath string) ([]models.Rider, error) {
	// Open the CSV file
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Create a new CSV reader
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	records = records[1:]
	riders := make([]models.Rider, 0, len(records))

	for _, record := range records {
		if len(record) != 7 {
			return nil, fmt.Errorf("invalid record length: expected 7, got %d", len(record))
		}
		dateOfBirth, err := time.Parse(time.RFC3339, record[5])
		if err != nil {
			return nil, fmt.Errorf("error parsing date of birth: %w", err)
		}
		createdAt, err := time.Parse(time.RFC3339, record[6])
		if err != nil {
			return nil, fmt.Errorf("error parsing created at: %w", err)
		}
		rider := models.Rider{
			ID:          record[0],
			FirstName:   record[1],
			LastName:    record[2],
			Email:       record[3],
			PhoneNumber: record[4],
			DateOfBirth: dateOfBirth,
			CreatedAt:   createdAt,
		}
		riders = append(riders, rider)
	}
	return riders, nil
}

func ExportDriversToCSV(drivers []models.Driver) {
	// Create CSV file for drivers
	driverFile, err := os.Create("data/drivers.csv")
	if err != nil {
		log.Fatal("Could not create drivers.csv file", err)
	}
	defer driverFile.Close()

	driverWriter := csv.NewWriter(driverFile)
	defer driverWriter.Flush()

	// Write header for drivers CSV
	driverWriter.Write([]string{"id", "first_name", "last_name", "email", "phone_number", "date_of_birth", "created_at"})

	// Generate and write drivers
	for _, driver := range drivers {
		driverWriter.Write([]string{
			driver.ID,
			driver.FirstName,
			driver.LastName,
			driver.Email,
			driver.PhoneNumber,
			driver.DateOfBirth.Format(time.RFC3339),
			driver.CreatedAt.Format(time.RFC3339),
		})
	}
}

func ImportDriversFromCSV(filePath string) ([]models.Driver, error) {
	// Open the CSV file
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Create a new CSV reader
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	records = records[1:]
	drivers := make([]models.Driver, 0, len(records))

	for _, record := range records {
		if len(record) != 7 {
			return nil, fmt.Errorf("invalid record length: expected 7, got %d", len(record))
		}
		dateOfBirth, err := time.Parse(time.RFC3339, record[5])
		if err != nil {
			return nil, fmt.Errorf("error parsing date of birth: %w", err)
		}
		createdAt, err := time.Parse(time.RFC3339, record[6])
		if err != nil {
			return nil, fmt.Errorf("error parsing created at: %w", err)
		}
		driver := models.Driver{
			ID:          record[0],
			FirstName:   record[1],
			LastName:    record[2],
			Email:       record[3],
			PhoneNumber: record[4],
			DateOfBirth: dateOfBirth,
			CreatedAt:   createdAt,
		}
		drivers = append(drivers, driver)
	}
	return drivers, nil
}
