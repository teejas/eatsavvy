package api

import (
	"fmt"
	"regexp"
)

func formatPhoneNumber(phone string) (string, error) {
	// Remove all non-digit characters
	re := regexp.MustCompile(`\D`)
	digits := re.ReplaceAllString(phone, "")

	// Handle 11-digit numbers starting with 1 (US country code)
	if len(digits) == 11 && digits[0] == '1' {
		digits = digits[1:]
	}

	// Validate we have exactly 10 digits
	if len(digits) != 10 {
		return "", fmt.Errorf("invalid phone number: expected 10 digits, got %d", len(digits))
	}

	// Format as (XXX) XXX-XXXX
	return fmt.Sprintf("(%s) %s-%s", digits[0:3], digits[3:6], digits[6:10]), nil
}
