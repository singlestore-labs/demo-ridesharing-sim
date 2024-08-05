package service

import "simulator/models"

func GenerateRiders(numRiders int) []models.Rider {
	riders := make([]models.Rider, numRiders)
	for i := 0; i < numRiders; i++ {
		riders[i] = models.GenerateRider()
	}
	return riders
}
