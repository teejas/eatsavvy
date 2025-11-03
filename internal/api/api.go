package api

import (
	"eatsavvy/pkg/db"
	"eatsavvy/pkg/places"
	"log/slog"
	netHttp "net/http"

	"github.com/gin-gonic/gin"
)

func StartServer(port string) {
	r := gin.Default()
	restaurantClient := places.NewRestaurantClient()
	dbClient := db.NewDatabaseClient()
	defer dbClient.Close()
	if dbClient == nil {
		slog.Error("[api.StartServer] Failed to create database client")
		return
	}
	r.POST("/search", func(c *gin.Context) {
		var request struct {
			Query string `json:"query"`
		}
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(netHttp.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		restaurants, err := restaurantClient.GetRestaurants(request.Query) // Magnin Cafe
		if err != nil {
			c.JSON(netHttp.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(netHttp.StatusOK, restaurants)
	})
	r.POST("/enrich", func(c *gin.Context) {
		var request struct {
			Id string `json:"id"`
		}
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(netHttp.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		err := restaurantClient.EnrichRestaurantDetails(request.Id)
		if err != nil {
			c.JSON(netHttp.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(netHttp.StatusOK, nil)
	})

	r.Run(":" + port)
}
