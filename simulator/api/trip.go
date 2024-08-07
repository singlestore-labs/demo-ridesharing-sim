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
		c.JSON(http.StatusOK, service.GetTripsByCityAndStatus(city, status))
	} else if city != "" {
		c.JSON(http.StatusOK, service.GetTripsByCity(city))
	} else if status != "" {
		c.JSON(http.StatusOK, service.GetTripsByStatus(status))
	} else {
		c.JSON(http.StatusOK, service.GetAllTrips())
	}
}
