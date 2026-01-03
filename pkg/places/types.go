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
	CookingOils           string `json:"oil"`
	NutFree               bool   `json:"nutFree"`
	DietaryAccommodations string `json:"accommodations"`
	Vegetables            string `json:"vegetables"`
}

type Restaurant struct {
	Id               string           `json:"id"`
	Name             string           `json:"name"`
	Address          string           `json:"address"`
	PhoneNumber      string           `json:"phoneNumber"`
	OpenHours        []TimeRange      `json:"openHours"`
	NutritionInfo    *NutritionInfo   `json:"nutritionInfo"`
	Rating           *float64         `json:"rating"`
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
	Rating              float64      `json:"rating"`
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

var GooglePlacesRestaurantTypes = []string{
	"acai_shop",
	"afghani_restaurant",
	"african_restaurant",
	"american_restaurant",
	"asian_restaurant",
	"bagel_shop",
	"bakery",
	"bar",
	"bar_and_grill",
	"barbecue_restaurant",
	"brazilian_restaurant",
	"breakfast_restaurant",
	"brunch_restaurant",
	"buffet_restaurant",
	"cafe",
	"cafeteria",
	"candy_store",
	"cat_cafe",
	"chinese_restaurant",
	"chocolate_factory",
	"chocolate_shop",
	"coffee_shop",
	"confectionery",
	"deli",
	"dessert_restaurant",
	"dessert_shop",
	"diner",
	"dog_cafe",
	"donut_shop",
	"fast_food_restaurant",
	"fine_dining_restaurant",
	"food_court",
	"french_restaurant",
	"greek_restaurant",
	"hamburger_restaurant",
	"ice_cream_shop",
	"indian_restaurant",
	"indonesian_restaurant",
	"italian_restaurant",
	"japanese_restaurant",
	"juice_shop",
	"korean_restaurant",
	"lebanese_restaurant",
	"meal_delivery",
	"meal_takeaway",
	"mediterranean_restaurant",
	"mexican_restaurant",
	"middle_eastern_restaurant",
	"pizza_restaurant",
	"pub",
	"ramen_restaurant",
	"restaurant",
	"sandwich_shop",
	"seafood_restaurant",
	"spanish_restaurant",
	"steak_house",
	"sushi_restaurant",
	"tea_house",
	"thai_restaurant",
	"turkish_restaurant",
	"vegan_restaurant",
	"vegetarian_restaurant",
	"vietnamese_restaurant",
	"wine_bar",
}
