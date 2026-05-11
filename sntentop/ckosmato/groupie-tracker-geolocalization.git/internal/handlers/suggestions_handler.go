package handlers

import (
	"encoding/json"
	"gtracker/internal/models"
	"gtracker/internal/services"
	"gtracker/utils"
	"log"
	"net/http"
	"slices"
	"strconv"
	"strings"
)

// SuggestionsHandler handles the request to fetch suggestions
func SuggestionsHandler(w http.ResponseWriter, r *http.Request) {
	query := utils.Replacer(strings.ToLower(r.URL.Query().Get("query")))
	allArtists, err := services.GetArtists()
	if err != nil {
		http.Error(w, "Failed to fetch artists", http.StatusInternalServerError)
		log.Printf("Error fetching artists: %v", err)
		return
	}

	suggestions := getSuggestions(query, allArtists, strings.HasPrefix)
	if len(suggestions) < 5 {
		suggestions = getSuggestions(query, allArtists, strings.Contains)
	}

	if suggestions == nil {
		suggestions = []string{} // Return an empty array instead of `null`
	}

	suggestions = Sorter(suggestions)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(struct {
		Suggestions []string `json:"suggestions"`
	}{Suggestions: suggestions})

	if err != nil {
		http.Error(w, `{"suggestions":[]}`, http.StatusInternalServerError) // Always return valid JSON
		log.Printf("Error encoding autocomplete response: %v", err)
	}
}

// Sorter sorts the suggestions
func Sorter(suggestions []string) []string {
	slices.Sort(suggestions)
	var cleanSugg []string
	for _, val := range suggestions {
		cleaned := strings.Split(val, "#")
		cleanSugg = append(cleanSugg, cleaned[1])
	}

	return cleanSugg
}

// matches checks if the query matches the value
func matches(query, value string, matchFunc func(string, string) bool) bool {
	return matchFunc(strings.ToLower(value), query)
}

// getSuggestions returns the suggestions for the query
func getSuggestions(query string, allArtists []models.Artist, matchfunc func(string, string) bool) []string {
	seenSuggestions := make(map[string]bool)
	var suggestions []string
	for _, artist := range allArtists {

		if len(suggestions) < 5 && matches(query, artist.Name, matchfunc) {
			if _, exists := seenSuggestions[artist.Name]; !exists {
				suggestions = append(suggestions, "1#"+artist.Name)
				seenSuggestions[artist.Name] = true
			}
		}

		for _, member := range artist.Members {
			if len(suggestions) < 5 && matches(query, member, matchfunc) {
				if _, exists := seenSuggestions[member]; !exists {
					suggestions = append(suggestions, "2#"+member)
					seenSuggestions[member] = true
				}
			}
		}
		if len(suggestions) < 5 && matches(query, utils.Replacer(artist.FirstAlbum), matchfunc) {
			if _, exists := seenSuggestions[artist.FirstAlbum]; !exists {
				suggestions = append(suggestions, "3#"+artist.FirstAlbum)
				seenSuggestions[artist.FirstAlbum] = true
			}

		}
		if len(suggestions) < 5 && matches(query, strconv.Itoa(artist.CreationDate), matchfunc) {
			if _, exists := seenSuggestions[strconv.Itoa(artist.CreationDate)]; !exists {
				suggestions = append(suggestions, "4#"+strconv.Itoa(artist.CreationDate))
				seenSuggestions[strconv.Itoa(artist.CreationDate)] = true
			}
		}

		for _, concert := range artist.ConcertsFormatted {
			concertSplit := strings.Split(concert, ":")

			if len(suggestions) < 5 && matches(query, concertSplit[0], matchfunc) {
				if _, exists := seenSuggestions[concertSplit[0]]; !exists {
					suggestions = append(suggestions, "5#"+concertSplit[0])
					seenSuggestions[concertSplit[0]] = true
				}

			}

		}

		if len(suggestions) == 5 {
			break
		}
	}
	return suggestions
}
