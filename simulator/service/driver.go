package service

import "simulator/models"

func GenerateDrivers(numDrivers int) []models.Driver {
	drivers := make([]models.Driver, numDrivers)
	for i := 0; i < numDrivers; i++ {
		drivers[i] = models.GenerateDriver()
	}
	return drivers
}
