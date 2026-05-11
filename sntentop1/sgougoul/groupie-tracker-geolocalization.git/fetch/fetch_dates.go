package fetch

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Constants for external API URLs
const datesURL = "https://groupietrackers.herokuapp.com/api/dates"

// Struct to match the dates API response
type Dates struct {
	ID    int      `json:"id"`
	Dates []string `json:"dates"`
}

// Fetch dates from the external API
func FetchDates(id int) (*Dates, error) {
	resp, err := http.Get(fmt.Sprintf("%s/%d", datesURL, id))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var dates Dates
	if err := json.NewDecoder(resp.Body).Decode(&dates); err != nil {
		return nil, err
	}
	return &dates, nil
}
