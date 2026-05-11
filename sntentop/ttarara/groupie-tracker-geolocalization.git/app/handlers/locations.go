package handlers

import (
	"encoding/json"
	"fmt"
	"groupie-tracker-geolocalization/app/services"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strings"
)

// HandleAllLocations returns a JSON array of *all* unique location strings
func HandleAllLocations(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Fetch the entire "locations" dataset from groupietrackers
	allLocIndex, err := services.FetchData[LocationsIndex]("https://groupietrackers.herokuapp.com/api/locations")
	if err != nil {
		http.Error(w, "Failed to fetch location data", http.StatusInternalServerError)
		return
	}

	// Build a set of unique location strings
	uniqueLocs := make(map[string]bool)
	for _, block := range allLocIndex.Index {
		for _, loc := range block.Locations {
			uniqueLocs[loc] = true
		}
	}

	// Convert set to slice
	var allLocs []string
	for locStr := range uniqueLocs {
		allLocs = append(allLocs, locStr)
	}

	// Sort them alphabetically:
	sort.Slice(allLocs, func(i, j int) bool {
		return strings.ToLower(allLocs[i]) < strings.ToLower(allLocs[j])
	})

	// Return as JSON array
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(allLocs)
}

// GeocodeAddress is the original function that performs geocoding.
func GeocodeAddress(address string) (float64, float64, error) {
	geocodeURL := fmt.Sprintf("https://maps.googleapis.com/maps/api/geocode/json?address=%s&key=AIzaSyC0ZG2NAoK8uT_lwyddrFSlZNEP_v2QkrA", url.QueryEscape(address))

	resp, err := http.Get(geocodeURL)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to geocode address: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to read geocoding response: %v", err)
	}

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
	if err := json.Unmarshal(body, &result); err != nil {
		return 0, 0, fmt.Errorf("failed to parse geocoding response: %v", err)
	}

	if result.Status != "OK" || len(result.Results) == 0 {
		return 0, 0, fmt.Errorf("no results found for the address")
	}

	return result.Results[0].Geometry.Location.Lat, result.Results[0].Geometry.Location.Lng, nil
}

// GeocodeAddressHandler is an HTTP handler that wraps GeocodeAddress.
func GeocodeAddressHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the address from the query parameters.
	address := r.URL.Query().Get("address")
	if address == "" {
		http.Error(w, "address query parameter is required", http.StatusBadRequest)
		return
	}

	// Call the GeocodeAddress function.
	lat, lng, err := GeocodeAddress(address)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to geocode address: %v", err), http.StatusInternalServerError)
		return
	}

	// Create a response struct.
	response := struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	}{
		Latitude:  lat,
		Longitude: lng,
	}

	// Marshal the response to JSON.
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, fmt.Sprintf("failed to encode response: %v", err), http.StatusInternalServerError)
		return
	}
}

// Construct the Google Maps API URL
func buildMapScriptURL() string {
	return fmt.Sprintf(
		"https://maps.googleapis.com/maps/api/js?key=%s&libraries=places,marker&loading=async&map_ids=500755a5e04d8f95",
		"AIzaSyC0ZG2NAoK8uT_lwyddrFSlZNEP_v2QkrA",
	)
}
