package handlers

import (
	"encoding/json"
	"gtracker/internal/services"
	"log"
	"net/http"
	"strings"
)

// FetchAvailableConcertLocations fetches the available concert locations
func FetchAvailableConcertLocations(w http.ResponseWriter, r *http.Request) {
	locationsSet := make(map[string]struct{})
	allArtists, err := services.GetArtists()
	if err != nil {
		log.Printf("Error fetching artists: %v", err)
		return
	}
	for _, artist := range allArtists {
		for _, concertLocation := range artist.ConcertsFormatted {
			parts := strings.Split(concertLocation, " : ")
			if len(parts) < 2 {
				continue
			}
			location := parts[0]
			locationsSet[location] = struct{}{}
		}
	}

	var locations []string
	for location := range locationsSet {
		locations = append(locations, location)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(locations); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
