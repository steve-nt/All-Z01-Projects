package backend

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// GeocodeAddress is a function that performs geocoding, converting an address into latitude and longitude coordinates.
func GeocodeAddress(address string) (float64, float64, error) {
	geocodeURL := fmt.Sprintf("https://maps.googleapis.com/maps/api/geocode/json?address=%s&key=AIzaSyDtuaYfzbNShwjWrDwBkEnhp2H3Jq9aG9g", url.QueryEscape(address))

	resp, err := http.Get(geocodeURL)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to geocode address: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to read geocoding response: %v", err)
	}

	// Define a struct to hold the geocoding API response.
	var result struct {
		Results []struct {
			Geometry struct {
				Location struct {
					Lat float64 `json:"lat"`
					Lng float64 `json:"lng"`
				} `json:"location"`
			} `json:"geometry"`
		} `json:"results"`
		Status string `json:"status"`
	}
	// Parse the JSON response into the defined struct.
	if err := json.Unmarshal(body, &result); err != nil {
		return 0, 0, fmt.Errorf("failed to parse geocoding response: %v", err)
	}

	if result.Status != "OK" || len(result.Results) == 0 {
		return 0, 0, fmt.Errorf("no results found for the address")
	}

	return result.Results[0].Geometry.Location.Lat, result.Results[0].Geometry.Location.Lng, nil
}

// GeocodeAddressHandler is an HTTP handler that wraps the GeocodeAddress function.
func GeocodeAddressHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the "address" query parameter from the URL.
	address := r.URL.Query().Get("address")
	if address == "" {
		http.Error(w, "address query parameter is required", http.StatusBadRequest)
		return
	}

	// Call the GeocodeAddress function to get the latitude and longitude.
	lat, lng, err := GeocodeAddress(address)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to geocode address: %v", err), http.StatusInternalServerError)
		return
	}

	// Define a struct to hold the response data.
	response := struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	}{
		Latitude:  lat,
		Longitude: lng,
	}

	// Marshal/Set the response to JSON.
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, fmt.Sprintf("failed to encode response: %v", err), http.StatusInternalServerError)
		return
	}
}
