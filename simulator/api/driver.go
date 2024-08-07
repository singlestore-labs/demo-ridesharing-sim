package api

import (
	"net/http"
	"simulator/service"

	"github.com/gin-gonic/gin"
)

func GetAllDrivers(c *gin.Context) {
	city := c.Query("city")
	if city != "" {
		c.JSON(http.StatusOK, service.GetDriversByCity(city))
	} else {
		c.JSON(http.StatusOK, service.GetAllDrivers())
	}
}
