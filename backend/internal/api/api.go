package api

import (
	"eatsavvy/internal/places"
	netHttp "net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

const MAX_ENRICHMENTS = 25

func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token != "Bearer "+os.Getenv("EATSAVVY_API_KEY") {
			c.JSON(netHttp.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
		}
		c.Next()
	}
}

func StartServer(port string) {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(gin.LoggerWithConfig(gin.LoggerConfig{
		SkipPaths: []string{"/health"},
	}))
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"https://eatsavvy.org", "http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))
	restaurantClient := places.NewRestaurantClient()
	defer restaurantClient.Close()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(netHttp.StatusOK, gin.H{"status": "ok"})
	})

	authorized := r.Group("/")
	authorized.Use(authMiddleware())

	authorized.GET("/restaurant/:id", func(c *gin.Context) {
		id := c.Param("id")
		restaurant, err := restaurantClient.GetRestaurant(id)
		if err != nil {
			c.JSON(netHttp.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(netHttp.StatusOK, restaurant)
	})

	authorized.PATCH("/restaurant/:id", func(c *gin.Context) {
		id := c.Param("id")
		var request struct {
			PhoneNumber string `json:"phoneNumber"`
		}
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(netHttp.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		formattedPhone, err := formatPhoneNumber(request.PhoneNumber)
		if err != nil {
			c.JSON(netHttp.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		restaurant, err := restaurantClient.UpdateRestaurantPhoneNumber(id, formattedPhone)
		if err != nil {
			c.JSON(netHttp.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(netHttp.StatusOK, restaurant)
	})

	authorized.GET("/restaurant", func(c *gin.Context) {
		restaurants, err := restaurantClient.GetAllRestaurants()
		if err != nil {
			c.JSON(netHttp.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(netHttp.StatusOK, restaurants)
	})

	authorized.POST("/search", func(c *gin.Context) {
		var request struct {
			Query string `json:"query"`
		}

		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(netHttp.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		restaurants, err := restaurantClient.SearchRestaurants(request.Query) // Magnin Cafe
		if err != nil {
			c.JSON(netHttp.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(netHttp.StatusOK, restaurants)
	})

	authorized.POST("/enrich", func(c *gin.Context) {
		var request struct {
			Ids []string `json:"ids"`
		}
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(netHttp.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		restaurants, err := restaurantClient.BatchEnrichRestaurantDetails(request.Ids)
		if err != nil {
			c.JSON(netHttp.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(netHttp.StatusOK, restaurants)
	})

	authorized.POST("/search-and-enrich", func(c *gin.Context) {
		var request struct {
			Query string `json:"query"`
		}
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(netHttp.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		restaurants, err := restaurantClient.SearchRestaurants(request.Query)
		if err != nil {
			c.JSON(netHttp.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if len(restaurants) == 0 {
			c.JSON(netHttp.StatusNotFound, gin.H{"error": "No restaurants found for query: " + request.Query})
			return
		}
		if len(restaurants) > MAX_ENRICHMENTS {
			c.JSON(netHttp.StatusBadRequest, gin.H{"error": "Too many restaurants found for query: " + request.Query + ". Please refine your query."})
			return
		}
		ids := []string{}
		for _, restaurant := range restaurants {
			ids = append(ids, restaurant.Id)
		}
		enrichedRestaurants, err := restaurantClient.BatchEnrichRestaurantDetails(ids)
		if err != nil {
			c.JSON(netHttp.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(netHttp.StatusOK, enrichedRestaurants)
	})

	authorized.POST("/process-eocr", func(c *gin.Context) {
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

	r.Run(":" + port)
}
