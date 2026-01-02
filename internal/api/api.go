package api

import (
	"eatsavvy/pkg/places"
	netHttp "net/http"

	"github.com/gin-gonic/gin"
)

func StartServer(port string) {
	r := gin.Default()
	restaurantClient := places.NewRestaurantClient()
	defer restaurantClient.Close()
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
		restaurant, err := restaurantClient.EnrichRestaurantDetails(request.Id)
		if err != nil {
			c.JSON(netHttp.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(netHttp.StatusOK, restaurant)
	})

	r.POST("/process-eocr", func(c *gin.Context) {
		var request places.EndOfCallReportMessage
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(netHttp.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		err := restaurantClient.UpdateRestaurantNutritionInfo(request)
		if err != nil {
			c.JSON(netHttp.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(netHttp.StatusOK, gin.H{"status": "ok"})
	})

	r.GET("/health", func(c *gin.Context) {
		c.JSON(netHttp.StatusOK, gin.H{"status": "ok"})
	})

	r.Run(":" + port)
}
