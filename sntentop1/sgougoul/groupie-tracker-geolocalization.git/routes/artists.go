package routes

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"sgougoupractice/fetch"
	"sgougoupractice/filters"
	"sgougoupractice/handlers"
)

// Fetch and serve the list of artists in JSON format
func serveArtists(w http.ResponseWriter, r *http.Request) error {
	artists, err := fetch.FetchArtists()
	if err != nil {
		log.Println("Error fetching artists:", err)
		return &handlers.HTTPError{
			Status:  http.StatusInternalServerError,
			Message: "error fetching artists",
		}
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(artists); err != nil {
		log.Println("Error encoding artists data:", err)
		return &handlers.HTTPError{
			Status:  http.StatusInternalServerError,
			Message: "bad encoding error",
		}
	}
	return nil
}

func serveFilterArtists(w http.ResponseWriter, r *http.Request) error {
	// Parse the filter criteria from the request body
	var filterData struct {
		CreationFrom     int    `json:"creationFrom"`
		CreationTo       int    `json:"creationTo"`
		AlbumFrom        int    `json:"albumFrom"`
		AlbumTo          int    `json:"albumTo"`
		SelectedMembers  []int  `json:"selectedMembers"`
		SelectedLocation string `json:"selectedLocation"`
	}

	if err := json.NewDecoder(r.Body).Decode(&filterData); err != nil {
		return &handlers.HTTPError{
			Status:  http.StatusBadRequest,
			Message: "Invalid filter data",
		}
	}

	// Fetch all artists
	artists, err := fetch.FetchArtists()
	if err != nil {
		return &handlers.HTTPError{
			Status:  http.StatusInternalServerError,
			Message: "Failed to fetch artists",
		}
	}

	// Convert location to slice
	var selectedLocations []string
	if filterData.SelectedLocation != "" {
		selectedLocations = []string{strings.TrimSpace(filterData.SelectedLocation)}
	}

	// Create filter options
	opts := filters.FilterOptions{
		CreationDateRange:   [2]int{filterData.CreationFrom, filterData.CreationTo},
		FirstAlbumYearRange: [2]int{filterData.AlbumFrom, filterData.AlbumTo},
		MemberCounts:        filterData.SelectedMembers,
		Locations:           selectedLocations,
	}

	// Use centralized filter logic
	filteredArtists := filters.FilterArtists(artists, opts)

	// Respond with filtered artists
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(filteredArtists)
	return nil
}
