package api

import (
	"net/http"
	"simulator/service"

	"github.com/gin-gonic/gin"
)

func GetAllDrivers(c *gin.Context) {
	c.JSON(http.StatusOK, service.GetAllDrivers())
}
