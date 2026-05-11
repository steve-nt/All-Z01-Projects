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

	var results []SearchResult

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
					results = append(results, SearchResult{
						Type:     "location",
						Display:  display,
						ArtistID: artist.ID,
					})
				}
			}
		}
	}

	// Search Artists and Members
	for _, artist := range artists {
		if strings.HasPrefix(strings.ToLower(artist.Name), searchQuery) {
			results = append(results, SearchResult{
				Type:     "artist",
				Display:  artist.Name + " - Artist/Band",
				ArtistID: artist.ID,
			})
		}

		// Search for members only if the band has more than one member and the artist name is different
		if len(artist.Members) > 1 {
			for _, member := range artist.Members {
				if strings.HasPrefix(strings.ToLower(member), searchQuery) && member != artist.Name {
					results = append(results, SearchResult{
						Type:     "member",
						Display:  fmt.Sprintf("%s - Member (%s)", member, artist.Name),
						ArtistID: artist.ID,
					})
				}
			}
		}

		// Search First Albums
		if strings.HasPrefix(strings.ToLower(artist.FirstAlbum), searchQuery) {
			results = append(results, SearchResult{
				Type:     "album",
				Display:  fmt.Sprintf("%s - First Album (%s)", artist.FirstAlbum, artist.Name),
				ArtistID: artist.ID,
			})
		}

		// Search Creation Date
		if strings.Contains(fmt.Sprint(artist.CreationDate), searchQuery) {
			results = append(results, SearchResult{
				Type:     "date",
				Display:  fmt.Sprintf("%d - Creation Date (%s)", artist.CreationDate, artist.Name),
				ArtistID: artist.ID,
			})
		}
	}

	// Return Results as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}
