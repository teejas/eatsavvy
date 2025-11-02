package places

import (
	"eatsavvy/pkg/http"
	"encoding/json"
	"log/slog"
	"os"
)

func GetPlacesDetails(textQuery string, fields []string) (Places, error) {
	httpClient := http.NewClient()

	reqBody := map[string]string{
		"textQuery": textQuery,
	}

	headers := map[string]string{
		"Content-Type":     "application/json",
		"X-Goog-Api-Key":   os.Getenv("GOOGLE_PLACES_API_KEY"),
		"X-Goog-FieldMask": GetGooglePlacesFieldMask(fields),
	}

	respBody, err := httpClient.Post("https://places.googleapis.com/v1/places:searchText", reqBody, headers)
	if err != nil {
		slog.Error("[places.GetPlacesDetails] Failed to send HTTP request", "error", err)
		return Places{}, err
	}

	var places Places
	err = json.Unmarshal(respBody, &places)
	if err != nil {
		slog.Error("[places.GetPlacesDetails] Failed to unmarshal response body", "error", err)
		return Places{}, err
	}

	return places, nil
}
