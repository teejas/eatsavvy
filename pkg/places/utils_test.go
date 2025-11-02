package places

import "testing"

func TestGetGooglePlacesFieldMask(t *testing.T) {
	fields := []string{"displayName", "currentOpeningHours", "currentSecondaryOpeningHours", "regularOpeningHours", "regularSecondaryOpeningHours", "nationalPhoneNumber", "restroom"}
	fieldMask := GetGooglePlacesFieldMask(fields)
	if fieldMask != "places.displayName,places.currentOpeningHours,places.currentSecondaryOpeningHours,places.regularOpeningHours,places.regularSecondaryOpeningHours,places.nationalPhoneNumber,places.restroom" {
		t.Errorf("Expected field mask to be 'places.displayName,places.currentOpeningHours,places.currentSecondaryOpeningHours,places.regularOpeningHours,places.regularSecondaryOpeningHours,places.nationalPhoneNumber,places.restroom', but got '%s'", fieldMask)
	}
}
