package main

import (
	"log/slog"

	"eatsavvy/pkg/places"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		slog.Error("Failed to load .env file", "error", err)
	}

	restaurants, err := places.GetRestaurantDetails("Magnin Cafe") // 138 Cyril Magnin St, San Francisco, CA 94102
	if err != nil {
		slog.Error("Failed to get restaurants", "error", err)
		return
	}

	slog.Info("Restaurants", "restaurants", restaurants)

}
