package api

import (
	"net/http"
	"simulator/service"

	"github.com/gin-gonic/gin"
)

func GetAllTrips(c *gin.Context) {
	city := c.Query("city")
	status := c.Query("status")
	if city != "" && status != "" {
	} else if city != "" {
	} else if status != "" {
		c.JSON(http.StatusOK, service.GetTripsByStatus(status))
	} else {
		c.JSON(http.StatusOK, service.GetAllTrips())
	}
}

func GetCurrentTripBreakdown(c *gin.Context) {
	tripsRequested := service.GetTripsByStatus("requested")
	tripsAccepted := service.GetTripsByStatus("accepted")
	tripsEnRoute := service.GetTripsByStatus("en_route")
	ridersIdle := service.GetRidersByStatus("idle")
	ridersRequested := service.GetRidersByStatus("requested")
	ridersWaiting := service.GetRidersByStatus("waiting")
	ridersInProress := service.GetRidersByStatus("in_progress")
	driversAvailable := service.GetDriversByStatus("available")
	driversInProress := service.GetDriversByStatus("in_progress")

	response := gin.H{
		"trips_requested":     len(tripsRequested),
		"trips_accepted":      len(tripsAccepted),
		"trips_en_route":      len(tripsEnRoute),
		"riders_idle":         len(ridersIdle),
		"riders_requested":    len(ridersRequested),
		"riders_waiting":      len(ridersWaiting),
		"riders_in_progress":  len(ridersInProress),
		"drivers_available":   len(driversAvailable),
		"drivers_in_progress": len(driversInProress),
	}
	c.JSON(http.StatusOK, response)
}
