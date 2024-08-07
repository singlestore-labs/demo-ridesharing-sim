package api

import (
	"net/http"
	"simulator/service"

	"github.com/gin-gonic/gin"
)

func GetAllTrips(c *gin.Context) {
	c.JSON(http.StatusOK, service.GetAllTrips())
}
