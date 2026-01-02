package places

import "time"

type EnrichmentStatus string

const (
	EnrichmentStatusPending    EnrichmentStatus = "pending"
	EnrichmentStatusInProgress EnrichmentStatus = "in_progress"
	EnrichmentStatusQueued     EnrichmentStatus = "queued"
	EnrichmentStatusCompleted  EnrichmentStatus = "completed"
	EnrichmentStatusFailed     EnrichmentStatus = "failed"
)

type NutritionInfo struct {
	CookingOils           string `json:"cookingOils"`
	NutAllergies          string `json:"nutAllergies"`
	DietaryAccommodations string `json:"dietaryAccommodations"`
	Vegetables            string `json:"vegetables"`
}

type Restaurant struct {
	Id               string           `json:"id"`
	Name             string           `json:"name"`
	Address          string           `json:"address"`
	PhoneNumber      string           `json:"phoneNumber"`
	OpenHours        []TimeRange      `json:"openHours"`
	NutritionInfo    []byte           `json:"nutritionInfo"`
	CreatedAt        time.Time        `json:"createdAt"`
	UpdatedAt        time.Time        `json:"updatedAt"`
	EnrichmentStatus EnrichmentStatus `json:"enrichmentStatus"`
}

type Places struct {
	Places []Place `json:"places"`
}

type Place struct {
	Id                  string       `json:"id"`
	PrimaryType         string       `json:"primaryType"`
	DisplayName         DisplayName  `json:"displayName"`
	Address             string       `json:"formattedAddress"`
	NationalPhoneNumber string       `json:"nationalPhoneNumber"`
	CurrentOpeningHours OpeningHours `json:"currentOpeningHours"`
	UtcOffsetMinutes    int          `json:"utcOffsetMinutes"`
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

// TimePoint represents a specific time on a weekday
type TimePoint struct {
	Weekday int `json:"weekday"`
	Hour    int `json:"hour"`
	Minute  int `json:"minute"`
}

// TimeRange represents an open/close time range
type TimeRange struct {
	Open  TimePoint `json:"open"`
	Close TimePoint `json:"close"`
}

type StructuredOutput struct {
	Name   string      `json:"name"`
	Result interface{} `json:"result"`
}

type EndOfCallReportMessage struct {
	Message struct {
		Artifact struct {
			Transcript        string                      `json:"transcript"`
			StructuredOutputs map[string]StructuredOutput `json:"structuredOutputs"`
		} `json:"artifact"`
		Analysis struct {
			Summary           string `json:summary`
			SuccessEvaluation string `json:"successEvaluation"`
		}
		Call struct {
			ID string `json:"id"`
		} `json:"call"`
		EndedReason string `json:"endedReason"`
	} `json:"message"`
}
