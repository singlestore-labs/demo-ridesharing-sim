package api

import (
	"net/http"
	"simulator/service"

	"github.com/gin-gonic/gin"
)

func GetAllRiders(c *gin.Context) {
	city := c.Query("city")
	if city != "" {
	} else {
		c.JSON(http.StatusOK, service.GetAllRiders())
	}
}
