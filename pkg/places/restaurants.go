package places

import (
	"log/slog"
)

func filterRestaurants(places []Place) []Place {
	restaurants := []Place{}
	for _, place := range places {
		if place.PrimaryType == "restaurant" {
			restaurants = append(restaurants, place)
		}
	}
	return restaurants
}

func GetRestaurantDetails(textQuery string) ([]Place, error) {
	fields := []string{
		"displayName",
		"primaryType",
		"currentOpeningHours",
		"regularOpeningHours",
		"nationalPhoneNumber",
	}
	places, err := GetPlacesDetails(textQuery, fields)
	if err != nil {
		slog.Error("[places.GetRestaurantDetails] Failed to get places details", "error", err)
		return nil, err
	}
	restaurants := filterRestaurants(places.Places)
	return restaurants, nil
}
