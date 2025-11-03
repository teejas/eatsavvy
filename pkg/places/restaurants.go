package places

import (
	"eatsavvy/pkg/db"
	"eatsavvy/pkg/http"
	"log/slog"
)

type RestaurantsClient struct {
	PlacesClient
	dbClient *db.DatabaseClient
}

func NewRestaurantClient() *RestaurantsClient {
	httpClient := http.NewClient()
	dbClient := db.NewDatabaseClient()
	return &RestaurantsClient{
		PlacesClient: PlacesClient{
			httpClient: httpClient,
		},
		dbClient: dbClient,
	}
}

func (rc *RestaurantsClient) GetRestaurants(textQuery string) ([]Place, error) {
	fields := []string{
		"id",
		"displayName",
		"primaryType",
	}
	places, err := rc.GetPlaces(textQuery, fields)
	if err != nil {
		slog.Error("[places.GetRestaurantDetails] Failed to get places details", "error", err)
		return nil, err
	}
	restaurants := filterRestaurants(places.Places)
	return restaurants, nil
}

func (rc *RestaurantsClient) EnrichRestaurantDetails(restaurantId string) error {
	restaurant, err := rc.getRestaurantDetails(restaurantId)
	slog.Info("[places.EnrichRestaurantDetails] Enriching details for restaurant", "restaurant", restaurant.DisplayName.Text)
	if err != nil {
		slog.Error("[places.EnrichRestaurantDetails] Failed to get restaurant details", "error", err)
		return err
	}
	_, err = rc.dbClient.Db.Exec(rc.dbClient.Ctx,
		"INSERT INTO public.restaurants (places_id, name, address, phone_number, open_hours) VALUES ($1, $2, $3, $4, $5) ON CONFLICT (places_id) DO NOTHING",
		restaurant.Id, restaurant.DisplayName.Text, restaurant.Address, restaurant.NationalPhoneNumber, restaurant.CurrentOpeningHours.Periods)
	if err != nil {
		slog.Error("[places.EnrichRestaurantDetails] Failed to insert restaurant details", "error", err)
		return err
	}
	// TODO: Enqueue job to enrich restaurant details
	return nil
}

func (rc *RestaurantsClient) getRestaurantDetails(restaurantId string) (Place, error) {
	fields := []string{
		"id",
		"displayName",
		"currentOpeningHours",
		"regularOpeningHours",
		"nationalPhoneNumber",
		"formattedAddress",
	}
	place, err := rc.GetPlaceDetails(restaurantId, fields)
	if err != nil {
		slog.Error("[places.GetRestaurantDetails] Failed to get place details", "error", err)
		return Place{}, err
	}
	return place, nil
}
