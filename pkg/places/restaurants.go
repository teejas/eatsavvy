package places

import (
	"eatsavvy/pkg/db"
	"eatsavvy/pkg/http"
	"eatsavvy/pkg/queue"
	"encoding/json"
	"time"

	"log/slog"
)

type RestaurantsClient struct {
	PlacesClient
	dbClient  *db.DatabaseClient
	publisher *queue.Publisher
}

func NewRestaurantClient() *RestaurantsClient {
	httpClient := http.NewClient()
	dbClient := db.NewDatabaseClient()
	publisher := queue.NewPublisher("enrich_restaurant_details")
	return &RestaurantsClient{
		PlacesClient: PlacesClient{
			httpClient: httpClient,
		},
		dbClient:  dbClient,
		publisher: publisher,
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
	err = rc.publisher.PublishMessage(restaurant)
	if err != nil {
		slog.Error("[places.EnrichRestaurantDetails] Failed to publish message", "error", err)
		return err
	}
	return nil
}

func (rc *RestaurantsClient) getRestaurantDetails(restaurantId string) (Place, error) {
	var place Place
	var name, address, phoneNumber string
	var openHours []byte
	var nutritionInfo []byte
	var createdAt, updatedAt time.Time

	err := rc.dbClient.Db.QueryRow(rc.dbClient.Ctx,
		"SELECT places_id, name, address, phone_number, open_hours, nutrition_info, created_at, updated_at FROM public.restaurants WHERE places_id = $1",
		restaurantId,
	).Scan(&place.Id, &name, &address, &phoneNumber, &openHours, &nutritionInfo, &createdAt, &updatedAt)

	if err == nil {
		place.DisplayName = DisplayName{Text: name}
		place.Address = address
		place.NationalPhoneNumber = phoneNumber
		if len(openHours) > 0 {
			json.Unmarshal(openHours, &place.CurrentOpeningHours.Periods)
		}
		if len(nutritionInfo) > 0 || updatedAt.After(time.Now().Add(-1*time.Hour*24*30)) { // 30 days
			return place, nil
		}
	}

	// If no row found, fetch from API
	fields := []string{
		"id",
		"displayName",
		"currentOpeningHours",
		"regularOpeningHours",
		"nationalPhoneNumber",
		"formattedAddress",
	}
	place, err = rc.GetPlaceDetails(restaurantId, fields)
	if err != nil {
		slog.Error("[places.GetRestaurantDetails] Failed to get place details", "error", err)
		return Place{}, err
	}

	_, err = rc.dbClient.Db.Exec(rc.dbClient.Ctx,
		`INSERT INTO public.restaurants (places_id, name, address, phone_number, open_hours) 
		VALUES ($1, $2, $3, $4, $5) 
		ON CONFLICT (places_id) DO UPDATE SET 
			name = EXCLUDED.name, 
			address = EXCLUDED.address, 
			phone_number = EXCLUDED.phone_number, 
			open_hours = EXCLUDED.open_hours,
			updated_at = NOW()`,
		place.Id, place.DisplayName.Text, place.Address, place.NationalPhoneNumber, place.CurrentOpeningHours.Periods)
	if err != nil {
		slog.Error("[places.EnrichRestaurantDetails] Failed to insert restaurant details", "error", err)
		return Place{}, err
	}
	return place, nil
}
