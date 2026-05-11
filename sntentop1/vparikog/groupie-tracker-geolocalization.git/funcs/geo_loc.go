package groupie_tracker_search

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// formatLocation replaces _ with space and - with comma
func formatLocation(raw string) string {
	return strings.ReplaceAll(strings.ReplaceAll(raw, "_", " "), "-", ", ")
}

// geocodeLocation queries Nominatim and returns lat/lon
func geocodeLocation(location string) (string, string, error) {
	baseURL := "https://nominatim.openstreetmap.org/search"
	params := url.Values{}
	params.Add("q", location)
	params.Add("format", "json")
	fullURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())

	client := http.Client{Timeout: 10 * time.Second}
	req, _ := http.NewRequest("GET", fullURL, nil)
	req.Header.Set("User-Agent", "GroupieTrackerApp/1.0")

	resp, err := client.Do(req)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	var result []GeoResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", "", err
	}

	if len(result) == 0 {
		return "", "", fmt.Errorf("no result found for: %s", location)
	}

	return result[0].Lat, result[0].Lon, nil
}

// GeolocateLocations processes all artist locations and returns coordinates
func GeolocateLocations(artists []Artist) []Artist {
	for i, artist := range artists {
		geoMap := make(map[string]GeoResult)

		for _, loc := range artist.Locations.Locations {
			formatted := formatLocation(loc)

			// Avoid duplicate geocoding if already processed
			if _, exists := geoMap[formatted]; exists {
				continue
			}

			lat, lon, err := geocodeLocation(formatted)
			if err != nil {
				fmt.Printf("Failed to geocode %s: %v\n", formatted, err)
				continue
			}

			geoMap[formatted] = GeoResult{Lat: lat, Lon: lon}
			fmt.Printf("✓ %s → (%s, %s)\n", formatted, lat, lon)
			time.Sleep(1 * time.Second) // Respect Nominatim rate limits
		}

		artists[i].GeoCoordinates = geoMap
	}
	return artists
}
