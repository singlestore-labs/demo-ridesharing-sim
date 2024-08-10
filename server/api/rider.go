package api

import (
	"server/service"

	"github.com/gin-gonic/gin"
)

func GetRiders(c *gin.Context) {
	c.JSON(200, service.GetAllRiders(c.Query("db")))
}
