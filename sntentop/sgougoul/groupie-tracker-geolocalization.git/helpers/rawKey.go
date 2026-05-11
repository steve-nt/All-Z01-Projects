package helpers

import (
	"strings"
)

// helper function for the geocoding worker for splitting city and country
func ParseKey(raw string) (city, country string) {
	parts := strings.Split(raw, "-")

	city = strings.ReplaceAll(parts[0], "_", " ")
	if len(parts) > 1 {
		country = parts[1]

	}

	return

}
