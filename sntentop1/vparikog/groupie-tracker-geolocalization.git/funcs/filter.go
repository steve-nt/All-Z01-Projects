package groupie_tracker_search

import (
	"sort"
	"strings"
)

// GetLocationMapping builds a map where each country is mapped to a slice of its unique cities.
func GetLocationMapping(s []Artist) map[string][]string {
	// We'll use a nested map to ensure uniqueness of cities for each country.
	countryToCities := make(map[string]map[string]bool)

	for _, record := range s {
		for _, loc := range record.Locations.Locations {
			parts := strings.Split(loc, "-")
			// Ensure we have at least two parts: city and country.
			if len(parts) < 2 {
				continue // Skip unexpected formats.
			}

			// Country: last part; City: join all parts before that.
			country := parts[len(parts)-1]
			city := strings.Join(parts[:len(parts)-1], "-")

			// Initialize the nested map if necessary.
			if countryToCities[country] == nil {
				countryToCities[country] = make(map[string]bool)
			}
			countryToCities[country][city] = true
		}
	}

	// Convert the nested maps to a map[string][]string for easier consumption.
	finalMapping := make(map[string][]string)
	for country, citiesMap := range countryToCities {
		for city := range citiesMap {
			finalMapping[country] = append(finalMapping[country], city)
		}
		// Sort the cities for each country.
		sort.Strings(finalMapping[country])
	}

	return finalMapping
}
