package backend

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

var (
	artists       Artists
	locationsList []string // Cache for locations
)

func HandleSearch(w http.ResponseWriter, r *http.Request) {
	searchQuery := strings.ToLower(r.URL.Query().Get("q"))

	if searchQuery == "" {
		w.Write([]byte("[]"))
		return
	}

	var results []string

	// Search Locations Directly from Artists
	locationResults := make(map[string]bool)
	for _, artist := range artists {
		for location := range artist.Relation.DatesLocations {
			formattedLocation := FormatLocation(location)
			if strings.Contains(strings.ToLower(formattedLocation), searchQuery) {
				// Show location with the artist's name
				display := fmt.Sprintf("%s (%s)", formattedLocation, artist.Name)
				if !locationResults[display] {
					locationResults[display] = true
					results = append(results, display)
				}
			}
		}
	}

	// Search Artists and Members
	for _, artist := range artists {
		if strings.Contains(strings.ToLower(artist.Name), searchQuery) {
			results = append(results, artist.Name+" - Artist/Band")
		}

		for _, member := range artist.Members {
			if strings.Contains(strings.ToLower(member), searchQuery) {
				results = append(results, fmt.Sprintf("%s - Member (%s)", member, artist.Name))
			}
		}

		// Search First Albums
		if strings.Contains(strings.ToLower(artist.FirstAlbum), searchQuery) {
			results = append(results, fmt.Sprintf("%s - First Album (%s)", artist.FirstAlbum, artist.Name))
		}

		// Search Creation Date
		if strings.Contains(fmt.Sprint(artist.CreationDate), searchQuery) {
			results = append(results, fmt.Sprintf("%d - Creation Date (%s)", artist.CreationDate, artist.Name))
		}
	}

	// Return Results as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}
