package api

import "github.com/gin-gonic/gin"

func GetRiders(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Hello, World!",
	})
}
