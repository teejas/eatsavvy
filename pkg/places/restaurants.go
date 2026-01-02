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

func (rc *RestaurantsClient) Close() {
	rc.dbClient.Close()
	rc.publisher.Close()
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

func (rc *RestaurantsClient) EnrichRestaurantDetails(restaurantId string) (Restaurant, error) {
	var restaurant Restaurant
	var openHours []byte
	var nutritionInfo []byte
	err := rc.dbClient.Db.QueryRow(rc.dbClient.Ctx,
		`SELECT places_id, name, address, phone_number, open_hours, nutrition_info, created_at, updated_at, enrichment_status 
		 FROM public.restaurants WHERE places_id = $1`,
		restaurantId,
	).Scan(&restaurant.Id, &restaurant.Name, &restaurant.Address, &restaurant.PhoneNumber, &openHours, &nutritionInfo, &restaurant.CreatedAt, &restaurant.UpdatedAt, &restaurant.EnrichmentStatus)
	if err == nil {
		if len(openHours) > 0 {
			json.Unmarshal(openHours, &restaurant.OpenHours)
		}
		// If enrichment is completed and updated within the last 30 days, or is in progress or queued, no need to enrich again
		if (restaurant.EnrichmentStatus == EnrichmentStatusCompleted &&
			restaurant.UpdatedAt.After(time.Now().Add(-1*time.Hour*24*30))) ||
			restaurant.EnrichmentStatus == EnrichmentStatusInProgress ||
			restaurant.EnrichmentStatus == EnrichmentStatusQueued {
			return restaurant, nil
		}
	}

	// If no row found, fetch from API
	fields := []string{
		"id",
		"displayName",
		"currentOpeningHours",
		"nationalPhoneNumber",
		"formattedAddress",
		"utcOffsetMinutes",
	}
	place, err := rc.GetPlaceDetails(restaurantId, fields)
	if err != nil {
		slog.Error("[places.EnrichRestaurantDetails] Failed to get place details", "error", err)
		return Restaurant{}, err
	}

	// Start transaction to check enrichment_status and upsert atomically
	tx, err := rc.dbClient.Db.Begin(rc.dbClient.Ctx)
	if err != nil {
		slog.Error("[places.EnrichRestaurantDetails] Failed to begin transaction", "error", err)
		return Restaurant{}, err
	}
	defer tx.Rollback(rc.dbClient.Ctx)

	// Check if enrichment_status is already "queued" or "in_progress"
	var existingStatus string
	err = tx.QueryRow(rc.dbClient.Ctx,
		`SELECT enrichment_status FROM public.restaurants WHERE places_id = $1 FOR UPDATE`,
		place.Id,
	).Scan(&existingStatus)

	if err == nil && (existingStatus == string(EnrichmentStatusQueued) || existingStatus == string(EnrichmentStatusInProgress)) {
		slog.Info("[places.EnrichRestaurantDetails] Skipping insert, enrichment already in progress or queued", "places_id", place.Id, "status", existingStatus)
		return restaurant, nil
	}

	restaurant.Id = place.Id
	restaurant.Name = place.DisplayName.Text
	restaurant.Address = place.Address
	if place.NationalPhoneNumber != "" {
		restaurant.PhoneNumber = place.NationalPhoneNumber
	}
	restaurant.OpenHours = periodsToTimeRanges(place.CurrentOpeningHours.Periods, place.UtcOffsetMinutes)
	restaurant.CreatedAt = time.Now()
	restaurant.UpdatedAt = time.Now()
	restaurant.EnrichmentStatus = EnrichmentStatusQueued

	// Proceed with upsert and set enrichment_status to "queued"
	_, err = tx.Exec(rc.dbClient.Ctx,
		`INSERT INTO public.restaurants (places_id, name, address, phone_number, open_hours, enrichment_status) 
		VALUES ($1, $2, $3, $4, $5, $6) 
		ON CONFLICT (places_id) DO UPDATE SET 
			name = EXCLUDED.name, 
			address = EXCLUDED.address, 
			phone_number = COALESCE(restaurants.phone_number, EXCLUDED.phone_number), 
			open_hours = EXCLUDED.open_hours,
			enrichment_status = EXCLUDED.enrichment_status,
			updated_at = NOW()
		`,
		place.Id, place.DisplayName.Text, place.Address, place.NationalPhoneNumber,
		periodsToTimeRanges(place.CurrentOpeningHours.Periods, place.UtcOffsetMinutes), EnrichmentStatusQueued,
	)
	if err != nil {
		slog.Error("[places.EnrichRestaurantDetails] Failed to insert restaurant details", "error", err)
		return Restaurant{}, err
	}

	if err = tx.Commit(rc.dbClient.Ctx); err != nil {
		slog.Error("[places.EnrichRestaurantDetails] Failed to commit transaction", "error", err)
		return Restaurant{}, err
	}

	slog.Info("[places.EnrichRestaurantDetails] Upserted restaurant", "places_id", restaurant.Id)

	err = rc.publisher.PublishMessage(restaurant)
	if err != nil {
		slog.Error("[places.EnrichRestaurantDetails] Failed to publish message", "error", err)
		_, err = rc.dbClient.Db.Exec(rc.dbClient.Ctx,
			`UPDATE public.restaurants SET enrichment_status = $1 WHERE places_id = $2`,
			EnrichmentStatusFailed, restaurant.Id,
		)
		if err != nil {
			slog.Error("[places.EnrichRestaurantDetails] Failed to update enrichment status", "error", err)
			return Restaurant{}, err
		}
		return Restaurant{}, err
	}

	slog.Info("[places.EnrichRestaurantDetails] Enqueued enrichment job for restaurant", "places_id", place.Id)

	return restaurant, nil
}

func (rc *RestaurantsClient) UpdateRestaurantNutritionInfo(eocr EndOfCallReportMessage) error {
	nutritionInfo := make(map[string]interface{})
	for _, result := range eocr.Message.Artifact.StructuredOutputs {
		nutritionInfo[result.Name] = result.Result
	}

	var placesId string

	err := rc.dbClient.Db.QueryRow(rc.dbClient.Ctx,
		`UPDATE public.calls SET call_status = $1, transcript = $2, structured_outputs = $3, summary = $4, success_evaluation = $5, ended_reason = $6, updated_at = NOW() WHERE vapi_call_id = $7 returning places_id`,
		"completed", eocr.Message.Artifact.Transcript, eocr.Message.Artifact.StructuredOutputs, eocr.Message.Analysis.Summary, eocr.Message.Analysis.SuccessEvaluation, eocr.Message.EndedReason, eocr.Message.Call.ID,
	).Scan(&placesId)
	if err != nil {
		return err
	}

	status := EnrichmentStatusCompleted
	if eocr.Message.Analysis.SuccessEvaluation == "false" || eocr.Message.EndedReason != "customer-ended-call" {
		status = EnrichmentStatusFailed
	}
	_, err = rc.dbClient.Db.Exec(rc.dbClient.Ctx,
		`UPDATE public.restaurants SET nutrition_info = $1, enrichment_status = $2, updated_at = NOW() WHERE places_id = $3`,
		nutritionInfo, status, placesId,
	)
	if err != nil {
		slog.Error("[places.UpdateRestaurantNutritionInfo] Failed to update restaurant nutrition info", "error", err)
		return err
	}
	slog.Info("[places.UpdateRestaurantNutritionInfo] Updated restaurant nutrition info", "places_id", placesId)
	return nil
}
