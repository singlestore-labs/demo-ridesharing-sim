package api

import (
	"net/http"
	"simulator/service"

	"github.com/gin-gonic/gin"
)

func GetAllRiders(c *gin.Context) {
	if c.Query("city") != "" {
		c.JSON(http.StatusOK, service.GetRidersInCity(c.Query("city")))
	} else {
		c.JSON(http.StatusOK, service.GetAllRiders())
	}
}
