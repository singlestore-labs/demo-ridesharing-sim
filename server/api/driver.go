package api

import (
	"server/service"

	"github.com/gin-gonic/gin"
)

func GetDrivers(c *gin.Context) {
	c.JSON(200, service.GetAllDrivers(c.Query("db")))
}
