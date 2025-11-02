package places

import "strings"

func GetGooglePlacesFieldMask(fields []string) string {
	for i := range fields {
		if !strings.HasPrefix(fields[i], "places.") {
			fields[i] = "places." + fields[i]
		}
	}
	return strings.Join(fields, ",")
}
