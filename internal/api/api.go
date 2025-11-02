package api

import (
	"eatsavvy/pkg/places"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Run() (*gin.Engine, error) {
	r := gin.Default()
	r.POST("/restaurants/search", func(c *gin.Context) {
		var request struct {
			Query string `json:"query"`
		}
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		restaurants, err := places.GetRestaurantDetails(request.Query)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, restaurants)
	})
	r.Run(":8080")

	return r, nil
}
