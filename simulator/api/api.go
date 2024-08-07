package api

import (
	"simulator/config"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var Router *gin.Engine

func StartServer() {
	r := gin.New()
	r.Use(
		gin.Recovery(),
		// gin.Logger(),
		cors.Default(),
	)
	RegisterRoutes(r)
	Router = r
	r.Run(":" + config.Port)
}

func RegisterRoutes(r *gin.Engine) {
	r.GET("/trips", GetAllTrips)
	r.GET("/trips/current", GetCurrentTripBreakdown)
	r.GET("/riders", GetAllRiders)
	r.GET("/drivers", GetAllDrivers)
}
