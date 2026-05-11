package fetch

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Constants for external API URLs
const relationURL = "https://groupietrackers.herokuapp.com/api/relation"

// Struct to match the relation API response
type Relation struct {
	ID             int                 `json:"id"`
	DatesLocations map[string][]string `json:"datesLocations"`
}

// Fetch relations from the external API
func FetchRelations(id int) (*Relation, error) {
	resp, err := http.Get(fmt.Sprintf("%s/%d", relationURL, id))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var relation Relation
	if err := json.NewDecoder(resp.Body).Decode(&relation); err != nil {
		return nil, err
	}
	return &relation, nil
}
