package geo

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
)

// var rateLimit = time.NewTicker(200 * time.Millisecond)

func Geocode(city, country string) (float64, float64, error) {
	// <-rateLimit.C

	apiKey := os.Getenv("GEO_API_KEY")
	if apiKey == "" {
		return 0, 0, fmt.Errorf("API key not found")
	}

	baseURL := os.Getenv("GEO_API_URL")
	if baseURL == "" {
		return 0, 0, fmt.Errorf("API URL not found")
	}
	params := url.Values{}
	params.Add("city", city)
	params.Add("country", country)
	fullURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())

	req, _ := http.NewRequest("GET", fullURL, nil)
	req.Header.Set("X-Api-Key", apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, 0, fmt.Errorf("API error: %s", resp.Status)
	}

	var results []struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		return 0, 0, err
	}

	if len(results) == 0 {
		return 0, 0, fmt.Errorf("no coordinates found")
	}

	return results[0].Latitude, results[0].Longitude, nil
}
