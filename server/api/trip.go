package api

import (
	"server/service"

	"github.com/gin-gonic/gin"
)

func GetCurrentTripStatus(c *gin.Context) {
	db := c.Query("db")
	city := c.Query("city")
	if city != "" {
		c.JSON(200, service.GetCurrentTripStatusByCity(db, city))
		return
	}
	c.JSON(200, service.GetCurrentTripStatus(db))
}
