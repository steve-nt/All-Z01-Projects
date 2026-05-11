package fetch

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type Relation struct {
	DatesLocations map[string][]string `json:"datesLocations"`
}

func fetchRelation(relationURL string) (map[string][]string, error) {
	resp, err := http.Get(relationURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch relation: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status code %d", resp.StatusCode)
	}

	var relData Relation
	if err := json.NewDecoder(resp.Body).Decode(&relData); err != nil {
		return nil, fmt.Errorf("failed to parse relation JSON: %w", err)
	}

	newDatesLocations := make(map[string][]string)

	for key, value := range relData.DatesLocations {
		newKey := strings.Replace(key, "_", " ", -1)
		newKey = strings.Replace(newKey, "-", ", ", -1)
		newDatesLocations[newKey] = value
	}

	relData.DatesLocations = newDatesLocations

	return relData.DatesLocations, nil
}
