package groupie_tracker_search

import (
	"encoding/json"
	"log"
	"net/http"
)

// Handler function to serve location mapping
func HandleLocations(w http.ResponseWriter, r *http.Request) {
	artists, err := getCachedArtists()

	if err != nil {
		http.Error(w, "Error fetching artists", http.StatusInternalServerError)
		log.Println("Error fetching artists:", err)
		return
	}
	// Assume `artists` is your dataset containing all artist data
	locations := GetLocationMapping(artists)

	// Convert to JSON and send response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(locations)
}
