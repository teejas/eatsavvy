package places

import (
	"eatsavvy/pkg/http"
	"encoding/json"
	"errors"
	"log/slog"
	"os"
)

type PlacesClient struct {
	httpClient *http.Http
}

func (pc *PlacesClient) GetPlaces(textQuery string, fields []string) (Places, error) {
	slog.Info("[places.GetPlaces] Getting places for text query", "textQuery", textQuery)
	reqBody := map[string]string{
		"textQuery": textQuery,
	}

	headers := map[string]string{
		"Content-Type":     "application/json",
		"X-Goog-Api-Key":   os.Getenv("GOOGLE_PLACES_API_KEY"),
		"X-Goog-FieldMask": getGooglePlacesFieldMask(fields, true),
	}

	respBody, statusCode, err := pc.httpClient.Post("https://places.googleapis.com/v1/places:searchText", reqBody, headers)
	if err != nil {
		slog.Error("[places.GetPlaces] Failed to send HTTP request", "error", err)
		return Places{}, err
	}
	if statusCode >= 400 {
		slog.Error("[places.GetPlaces] Failed to get places", "statusCode", statusCode, "responseBody", string(respBody))
		return Places{}, errors.New("failed to get places: " + string(respBody))
	}

	var places Places
	err = json.Unmarshal(respBody, &places)
	if err != nil {
		slog.Error("[places.GetPlaces] Failed to unmarshal response body", "error", err)
		return Places{}, err
	}

	return places, nil
}

func (pc *PlacesClient) GetPlaceDetails(placeId string, fields []string) (Place, error) {
	slog.Info("[places.GetPlaceDetails] Getting place details for place", "placeId", placeId)

	headers := map[string]string{
		"Content-Type":     "application/json",
		"X-Goog-Api-Key":   os.Getenv("GOOGLE_PLACES_API_KEY"),
		"X-Goog-FieldMask": getGooglePlacesFieldMask(fields, false),
	}

	respBody, statusCode, err := pc.httpClient.Get("https://places.googleapis.com/v1/places/"+placeId, headers)
	if err != nil {
		slog.Error("[places.GetPlaceDetails] Failed to send HTTP request", "error", err)
		return Place{}, err
	}
	if statusCode >= 400 {
		slog.Error("[places.GetPlaceDetails] Failed to get place details", "statusCode", statusCode, "responseBody", string(respBody))
		return Place{}, errors.New("failed to get place details: " + string(respBody))
	}

	var place Place
	err = json.Unmarshal(respBody, &place)
	if err != nil {
		slog.Error("[places.GetPlaceDetails] Failed to unmarshal response body", "error", err)
		return Place{}, err
	}

	return place, nil
}
