package api

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		MaxAge:           12 * time.Hour,
		AllowCredentials: true,
	}))
	InitializeRoutes(r)
	return r
}

func InitializeRoutes(router *gin.Engine) {
	router.GET("/trips/current/status", GetCurrentTripStatus)
	router.GET("/trips/statistics", GetTotalTripStatistics)
	router.GET("/trips/statistics/daily", GetDailyTripStatistics)
	router.GET("/trips/last/hour", GetMinuteTripCountsLastHour)
	router.GET("/trips/last/day", GetHourlyTripCountsLastDay)
	router.GET("/trips/last/week", GetDailyTripCountsLastWeek)
	router.GET("/riders", GetRiders)
	router.GET("/drivers", GetDrivers)
	router.GET("/cities", GetCities)
}
