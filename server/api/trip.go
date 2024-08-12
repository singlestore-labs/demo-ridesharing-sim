package api

import (
	"server/service"

	"github.com/gin-gonic/gin"
)

func GetCurrentTripStatus(c *gin.Context) {
	city := c.Query("city")
	db := c.Query("db")
	if city != "" {
		c.JSON(200, service.GetCurrentTripStatusByCity(city, db))
	}
	c.JSON(200, service.GetCurrentTripStatus(db))
}
