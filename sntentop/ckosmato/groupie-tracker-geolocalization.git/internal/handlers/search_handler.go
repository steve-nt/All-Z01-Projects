package handlers

import (
	"gtracker/internal/models"
	"gtracker/internal/services"
	"gtracker/utils"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// SearchHandler handles the request to search for artists
func SearchHandler(w http.ResponseWriter, r *http.Request) {
	query := utils.Replacer(strings.ToLower(r.URL.Query().Get("search")))

	var matchedArtist []models.Artist
	allArtists, err := services.GetArtists()
	if err != nil {
		http.Error(w, "Failed to fetch artists", http.StatusInternalServerError)
		log.Printf("Error fetching artists: %v", err)
		return
	}

	seenMatch := make(map[string]bool)

	for _, artist := range allArtists {
		if strings.Contains(strings.ToLower(artist.Name), query) && !seenMatch[artist.Name] {
			matchedArtist = append(matchedArtist, artist)
			seenMatch[artist.Name] = true

		}
		for i := range artist.Members {
			if strings.Contains(strings.ToLower(artist.Members[i]), query) && !seenMatch[artist.Name] {
				matchedArtist = append(matchedArtist, artist)
				seenMatch[artist.Name] = true

			}
		}
		if strings.HasPrefix(utils.Replacer(artist.FirstAlbum), query) && !seenMatch[artist.Name] {
			matchedArtist = append(matchedArtist, artist)
			seenMatch[artist.Name] = true
		}

		if strings.HasPrefix(strconv.Itoa(artist.CreationDate), query) && !seenMatch[artist.Name] {
			matchedArtist = append(matchedArtist, artist)
			seenMatch[artist.Name] = true

		}
		for i := range artist.ConcertsFormatted {
			if strings.Contains(strings.ToLower(artist.ConcertsFormatted[i]), query) && !seenMatch[artist.Name] {
				matchedArtist = append(matchedArtist, artist)
				seenMatch[artist.Name] = true
			}
		}
	}

	if len(matchedArtist) == 0 {

		if err := services.NoResultsTemplate.Execute(w, nil); err != nil {
			http.Error(w, "Failed to render no results page", http.StatusInternalServerError)
			log.Printf("Error executing noresults template: %v", err)
		}
		return
	}

	if err := services.IndexTemplate.Execute(w, matchedArtist); err != nil {
		http.Error(w, "Failed to render page", http.StatusInternalServerError)
		log.Printf("Error executing template: %v", err)
	}
}
