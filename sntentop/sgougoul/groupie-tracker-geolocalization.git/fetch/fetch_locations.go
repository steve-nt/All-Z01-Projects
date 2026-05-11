package fetch

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Constants for external API URLs
const locationsURL = "https://groupietrackers.herokuapp.com/api/locations"

// Struct to match the location API response
type Location struct {
	ID        int      `json:"id"`
	Locations []string `json:"locations"` // List of locations
}

// Fetch location data based on artist ID
func FetchLocations(id int) (*Location, error) {
	// Make the API request to fetch locations
	resp, err := http.Get(fmt.Sprintf("%s/%d", locationsURL, id))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Decode the JSON response into the Location struct
	var location Location
	if err := json.NewDecoder(resp.Body).Decode(&location); err != nil {
		return nil, err
	}

	return &location, nil
}
