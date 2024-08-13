package api

import (
	"server/service"

	"github.com/gin-gonic/gin"
)

func GetCities(c *gin.Context) {
	db := c.Query("db")
	c.JSON(200, service.GetCities(db))
}

func GetCurrentTripStatus(c *gin.Context) {
	db := c.Query("db")
	city := c.Query("city")
	if city != "" {
		c.JSON(200, service.GetCurrentTripStatusByCity(db, city))
		return
	}
	c.JSON(200, service.GetCurrentTripStatus(db))
}

func GetMinuteTripCountsLastHour(c *gin.Context) {
	db := c.Query("db")
	city := c.Query("city")
	c.JSON(200, service.GetMinuteTripCountsLastHour(db, city))
}

func GetHourlyTripCountsLastDay(c *gin.Context) {
	db := c.Query("db")
	city := c.Query("city")
	c.JSON(200, service.GetHourlyTripCountsLastDay(db, city))
}
