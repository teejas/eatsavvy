package places

type Places struct {
	Places []Place `json:"places"`
}

type Place struct {
	PrimaryType         string       `json:"primaryType"`
	DisplayName         DisplayName  `json:"displayName"`
	NationalPhoneNumber string       `json:"nationalPhoneNumber"`
	CurrentOpeningHours OpeningHours `json:"currentOpeningHours"`
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
