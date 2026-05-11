package groupie_tracker_search

import (
	"encoding/json"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

func SearchHandler(w http.ResponseWriter, r *http.Request) {
	query := strings.ToLower(r.URL.Query().Get("q"))

	if query == "" {
		http.Error(w, "Query cannot be empty", http.StatusBadRequest)
		return
	}

	var results []map[string]interface{} // Change to include interface{} to hold both strings and integers

	for _, artist := range cachedArtists {
		// Search for artist name
		if strings.Contains(strings.ToLower(artist.Name), query) {
			results = append(results, map[string]interface{}{
				"id":   artist.ID, // Add the ID
				"name": artist.Name,
				"type": "artist/band",
				"Prio": 100,
			})
		}

		// Search for members
		for _, member := range artist.Members {
			if strings.Contains(strings.ToLower(member), query) {
				if member == artist.Name && len(artist.Members) == 1 {
					continue
				}
				results = append(results, map[string]interface{}{
					"id":   artist.ID, // Link to artist ID
					"name": member,
					"type": "member - " + artist.Name,
					"Prio": 90,
				})
			}
		}

		// Search for locations
		for _, location := range artist.Locations.Locations {
			if strings.Contains(strings.ToLower(location), query) {
				results = append(results, map[string]interface{}{
					"id":   artist.ID, // Link to artist ID
					"name": location,
					"type": "location - " + artist.Name,
					"Prio": 80,
				})
			}
		}

		// Search for dates
		for _, date := range artist.Dates.Dates {
			if strings.Contains(strings.ToLower(date), query) {
				results = append(results, map[string]interface{}{
					"id":   artist.ID, // Link to artist ID
					"name": date,
					"type": "concert's date - " + artist.Name,
					"Prio": 70,
				})
			}
		}

		// Search for first album date
		if strings.Contains(strings.ToLower(artist.FirstAlbum), query) {
			results = append(results, map[string]interface{}{
				"id":   artist.ID, // Link to artist ID
				"name": artist.FirstAlbum,
				"type": "first album date - " + artist.Name,
				"Prio": 60,
			})
		}

		// Search for creation date
		if strings.Contains(strings.ToLower(strconv.Itoa(artist.CreationDate)), query) {
			results = append(results, map[string]interface{}{
				"id":   artist.ID, // Link to artist ID
				"name": strconv.Itoa(artist.CreationDate),
				"type": "Creation Date: " + artist.Name,
				"Prio": 50,
			})

		}
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i]["Prio"].(int) > results[j]["Prio"].(int)
	})

	// If no results found
	if len(results) == 0 {
		http.Error(w, "No results found", http.StatusNotFound)
		return
	}

	// Return results as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}
