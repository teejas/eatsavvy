package places

import "strings"

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
