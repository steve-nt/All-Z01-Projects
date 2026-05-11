package data

import (
	"encoding/json"
	"groupie-tracker/server/models"
	"net/http"
)

func FetchAllData() (models.Indexes, error) {
	var indexes models.Indexes

	// Fetch Artists - handle array response directly
	artistsURL := "https://groupietrackers.herokuapp.com/api/artists"
	resp, err := http.Get(artistsURL)
	if err != nil {
		return indexes, err
	}
	defer resp.Body.Close()

	// Decode the JSON array directly into Artists slice
	if err := json.NewDecoder(resp.Body).Decode(&indexes.Artists); err != nil {
		return indexes, err
	}

	// Fetch Locations
	locationsURL := "https://groupietrackers.herokuapp.com/api/locations"
	resp, err = http.Get(locationsURL)
	if err != nil {
		return indexes, err
	}
	defer resp.Body.Close()

	var locationData models.LocationData
	if err := json.NewDecoder(resp.Body).Decode(&locationData); err != nil {
		return indexes, err
	}
	indexes.Locations = locationData.Index

	// Fetch Dates
	datesURL := "https://groupietrackers.herokuapp.com/api/dates"
	resp, err = http.Get(datesURL)
	if err != nil {
		return indexes, err
	}
	defer resp.Body.Close()

	var dateData models.DateData
	if err := json.NewDecoder(resp.Body).Decode(&dateData); err != nil {
		return indexes, err
	}
	indexes.Dates = dateData.Index

	// Fetch Relations
	relationsURL := "https://groupietrackers.herokuapp.com/api/relation"
	resp, err = http.Get(relationsURL)
	if err != nil {
		return indexes, err
	}
	defer resp.Body.Close()

	var relationData models.RelationData
	if err := json.NewDecoder(resp.Body).Decode(&relationData); err != nil {
		return indexes, err
	}
	indexes.Relations = relationData.Index

	return indexes, nil
}
