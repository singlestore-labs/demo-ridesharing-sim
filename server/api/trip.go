package api

import (
	"server/service"

	"github.com/gin-gonic/gin"
)

func GetCurrentTripStatus(c *gin.Context) {
	c.JSON(200, service.GetCurrentTripStatus(c.Query("db")))
}
