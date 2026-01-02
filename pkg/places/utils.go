package places

import (
	"sort"
	"strings"
)

func getGooglePlacesFieldMask(fields []string, needsPlacesPrefix bool) string {
	if needsPlacesPrefix { // If using the Text Search API, we need to prefix the fields with "places."
		for i := range fields {
			if !strings.HasPrefix(fields[i], "places.") {
				fields[i] = "places." + fields[i]
			}
		}
	}
	return strings.Join(fields, ",")
}

func filterRestaurants(places []Place) []Place {
	restaurants := []Place{}
	for _, place := range places {
		if strings.Contains(place.PrimaryType, "restaurant") {
			restaurants = append(restaurants, place)
		}
	}
	return restaurants
}

// applyUtcOffset converts a local TimePoint to UTC by subtracting the UTC offset
func applyUtcOffset(tp TimePoint, utcOffsetMinutes int) TimePoint {
	// Convert to total minutes, subtract offset to get UTC
	totalMinutes := tp.Hour*60 + tp.Minute - utcOffsetMinutes

	// Handle day wraparound
	dayOffset := 0
	if totalMinutes < 0 {
		dayOffset = -1
		totalMinutes += 24 * 60
	} else if totalMinutes >= 24*60 {
		dayOffset = 1
		totalMinutes -= 24 * 60
	}

	newWeekday := (tp.Weekday + dayOffset + 7) % 7

	return TimePoint{
		Weekday: newWeekday,
		Hour:    totalMinutes / 60,
		Minute:  totalMinutes % 60,
	}
}

// periodsToTimeRanges converts a slice of Period to a slice of TimeRange in UTC, sorted by weekday index
func periodsToTimeRanges(periods []Period, utcOffsetMinutes int) []TimeRange {
	ranges := make([]TimeRange, len(periods))
	for i, p := range periods {
		open := TimePoint{Weekday: p.Open.Day, Hour: p.Open.Hour, Minute: p.Open.Minute}
		close := TimePoint{Weekday: p.Close.Day, Hour: p.Close.Hour, Minute: p.Close.Minute}

		ranges[i] = TimeRange{
			Open:  applyUtcOffset(open, utcOffsetMinutes),
			Close: applyUtcOffset(close, utcOffsetMinutes),
		}
	}
	// Sort by weekday index
	sort.Slice(ranges, func(i, j int) bool {
		return ranges[i].Open.Weekday < ranges[j].Open.Weekday
	})
	return ranges
}
