package places

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"os"
)

type RestaurantDetails struct {
	Places []Place `json:"places"`
}

type Place struct {
	CurrentOpeningHours OpeningHours `json:"currentOpeningHours"`
	DisplayName         DisplayName  `json:"displayName"`
	NationalPhoneNumber string       `json:"nationalPhoneNumber"`
	RegularOpeningHours OpeningHours `json:"regularOpeningHours"`
}

type OpeningHours struct {
	NextOpenTime        string   `json:"nextOpenTime"`
	OpenNow             bool     `json:"openNow"`
	Periods             []Period `json:"periods"`
	WeekdayDescriptions []string `json:"weekdayDescriptions"`
}

type Period struct {
	Open  TimeSlot `json:"open"`
	Close TimeSlot `json:"close"`
}

type TimeSlot struct {
	Date   *Date `json:"date,omitempty"`
	Day    int   `json:"day"`
	Hour   int   `json:"hour"`
	Minute int   `json:"minute"`
}

type Date struct {
	Day   int `json:"day"`
	Month int `json:"month"`
	Year  int `json:"year"`
}

type DisplayName struct {
	LanguageCode string `json:"languageCode"`
	Text         string `json:"text"`
}

func GetRestaurantDetails(textQuery string) []Place {
	client := &http.Client{}

	reqBody := map[string]string{
		"textQuery": textQuery,
	}

	var googlePlacesFields = []string{
		"displayName",
		"currentOpeningHours",
		"currentSecondaryOpeningHours",
		"regularOpeningHours",
		"regularSecondaryOpeningHours",
		"nationalPhoneNumber",
	}

	jsonReqBody, err := json.Marshal(reqBody)
	if err != nil {
		slog.Error("[GetRestaurantDetails] Failed to marshal request body", "error", err)
	}

	req, err := http.NewRequest("POST", "https://places.googleapis.com/v1/places:searchText", bytes.NewBuffer(jsonReqBody))
	if err != nil {
		slog.Error("[GetRestaurantDetails] Failed to create HTTP request", "error", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Goog-Api-Key", os.Getenv("GOOGLE_PLACES_API_KEY"))
	req.Header.Set("X-Goog-FieldMask", GetGooglePlacesFieldMask(googlePlacesFields))

	resp, err := client.Do(req)
	if err != nil {
		slog.Error("[GetRestaurantDetails] Failed to send HTTP request", "error", err)
	}
	defer resp.Body.Close()

	slog.Info("[GetRestaurantDetails] Response status", "status", resp.Status)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("[GetRestaurantDetails] Failed to read response body", "error", err)
	}

	slog.Info("[GetRestaurantDetails] Response body", "body", string(body))

	var response RestaurantDetails
	err = json.Unmarshal(body, &response)
	if err != nil {
		slog.Error("[GetRestaurantDetails] Failed to unmarshal response body", "error", err)
	}

	return response.Places
}
