package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"
	"time"
)

const mapTilerAPIKey = "zQeY1QnryK4tr9gEl0Ht"

type GeocodeResponse struct {
	Features []struct {
		Geometry struct {
			Coordinates []float64 `json:"coordinates"`
		} `json:"geometry"`
	} `json:"features"`
}

type CachedLocation struct {
	Lat string
	Lon string
}

var (
	geoCache = make(map[string]CachedLocation)
	cacheMux sync.RWMutex
)

// Function to get latitude and longitude from MapTiler with caching
func GeocodeAddress(location string) (string, string, error) {

	cacheMux.RLock()
	if coords, found := geoCache[location]; found {
		cacheMux.RUnlock()
		return coords.Lat, coords.Lon, nil
	}
	cacheMux.RUnlock()

	encodedLocation := url.QueryEscape(location)
	fullURL := fmt.Sprintf("https://api.maptiler.com/geocoding/%s.json?key=%s", encodedLocation, mapTilerAPIKey)

	client := &http.Client{
		Timeout: 1 * time.Second,
	}
	resp, err := client.Get(fullURL)
	if err != nil {
		return "", "", fmt.Errorf("failed to fetch geolocation: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", fmt.Errorf("failed to read response: %v", err)
	}

	var result GeocodeResponse
	err = json.Unmarshal(body, &result)
	if err != nil || len(result.Features) == 0 {
		return "", "", fmt.Errorf("invalid geolocation response")
	}

	lon := fmt.Sprintf("%f", result.Features[0].Geometry.Coordinates[0])
	lat := fmt.Sprintf("%f", result.Features[0].Geometry.Coordinates[1])

	cacheMux.Lock()
	geoCache[location] = CachedLocation{Lat: lat, Lon: lon}
	cacheMux.Unlock()

	return lat, lon, nil
}
