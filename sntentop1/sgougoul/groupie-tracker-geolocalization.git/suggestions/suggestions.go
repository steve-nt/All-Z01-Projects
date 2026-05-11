package suggestions

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"

	"sgougoupractice/fetch"
	"sgougoupractice/helpers"
)

// Suggestions represents the data structure for an artist or band suggestion, including name, members, and other attributes.
type Suggestions struct {
	Name         string   `json:"name"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
	Locations    string   `json:"locations"`
	ConcertDates string   `json:"-"`
}

type TypeSuggestion struct {
	Label string `json:"label"`
	Type  string `json:"type"`
}

// Global variables for caching suggestions and managing concurrency
var (
	Cache            []fetch.Artist
	CacheLock        sync.RWMutex
	SuggestionsCache []TypeSuggestion
)

// InitCache initializes the suggestions cache by fetching artist data and enriching the locations.
func InitCache() {
	// Fetch the list of artists from the data source
	data, err := fetch.FetchArtists()
	if err != nil {
		log.Println("Error fetching artists during initialization:", err)
		return
	}
	// Lock the cache and update the artist data
	CacheLock.Lock()
	Cache = data
	CacheLock.Unlock()
	log.Println("Artist suggestions cache initialized.")
	// Enrich the fetched artist data with location information
	EnrichLocations()
}

// EnrichLocations enriches the artist data by fetching the locations for each artist asynchronously.
func EnrichLocations() {
	// Lock the cache for reading
	CacheLock.RLock()
	artists := make([]fetch.Artist, len(Cache))
	copy(artists, Cache)
	CacheLock.RUnlock()

	var wg sync.WaitGroup

	updatedArtists := make([]fetch.Artist, len(artists))

	var mu sync.Mutex
	// Fetch locations asynchronously for each artist
	for i, artist := range artists {
		wg.Add(1)
		go func(i int, artist fetch.Artist) {
			defer wg.Done()
			// Fetch location data for the artist
			locationData, err := fetch.FetchLocations(artist.ID)
			if err != nil {

			} else if locationData != nil && len(locationData.Locations) > 0 {

				artist.Locations = strings.Join(locationData.Locations, ",")
				//log.Printf("Fetched locations for artist %d: %v", artist.ID, locationData.Locations)
			} else {
				artist.Locations = "Unavailable"
			}
			// Lock to update the updatedArtists slice
			mu.Lock()
			updatedArtists[i] = artist
			mu.Unlock()
		}(i, artist)
	}
	// Wait for all goroutines to finish fetching locations
	wg.Wait()
	// Lock the cache for updating
	CacheLock.Lock()
	Cache = updatedArtists
	// Build suggestions from the enriched artist data
	SuggestionsCache = BuildAllSuggestions(updatedArtists)
	CacheLock.Unlock()
	log.Println("Artist cache enriched with actual locations data.")
}

// RefreshCache refreshes the cache by periodically fetching new artist data.
func RefreshCache(interval time.Duration) {
	for {
		time.Sleep(interval)
		data, err := fetch.FetchArtists()
		if err != nil {
			log.Println("Error refreshing suggestions cache:", err)
			continue
		}

		CacheLock.RLock()
		Cache = data
		CacheLock.RUnlock()
		log.Println("Suggestions cache refreshed.")
	}

}

// BuildSuggestionsTypes constructs a list of suggestion types for a given artist.
func BuildSuggestionsTypes(artist fetch.Artist) []TypeSuggestion {
	var suggestions []TypeSuggestion
	// Add the artist/band itself to suggestions
	suggestions = append(suggestions, TypeSuggestion{
		Label: artist.Name,
		Type:  "artist/band",
	})
	// Add each band member to suggestions
	for _, member := range artist.Members {
		suggestions = append(suggestions, TypeSuggestion{
			Label: member,
			Type:  "member",
		})
	}
	// Add the artist's creation date to suggestions (if available)
	if artist.CreationDate != 0 {
		suggestions = append(suggestions, TypeSuggestion{
			Label: strconv.Itoa(artist.CreationDate),
			Type:  "CreationDate",
		})
	}
	// Add the first album to suggestions (if available)
	if trimmed := strings.TrimSpace(artist.FirstAlbum); trimmed != "" {
		log.Printf("%v", trimmed)
		suggestions = append(suggestions, TypeSuggestion{
			Label: trimmed,
			Type:  "firstAlbum",
		})
	}
	// Add the artist's locations to suggestions (if available)
	if trimmed := strings.TrimSpace(artist.Locations); trimmed != "" {
		locs := strings.Split(trimmed, ",")
		for _, loc := range locs {
			loc = strings.TrimSpace(loc)
			if loc != "" {
				suggestions = append(suggestions, TypeSuggestion{
					Label: loc,
					Type:  "location",
				})
			}
		}

	}

	return suggestions
}

// BuildAllSuggestions constructs a list of unique suggestions from the cached artist data.
func BuildAllSuggestions(cachedArtist []fetch.Artist) []TypeSuggestion {
	var suggestions1 []TypeSuggestion
	uniqueMap := make(map[string]TypeSuggestion)
	// Collect suggestions for each artist
	for _, artist := range cachedArtist {
		suggestions1 = append(suggestions1, BuildSuggestionsTypes(artist)...)
	}
	// Remove duplicates by using a map
	for _, suggestion := range suggestions1 {
		key := fmt.Sprintf("%v", suggestion)
		if _, exists := uniqueMap[key]; !exists {
			uniqueMap[key] = suggestion
		}
	}
	// Convert the map of unique suggestions into a slice
	suggestions := make([]TypeSuggestion, 0, len(uniqueMap))
	for _, suggestion := range uniqueMap {
		suggestions = append(suggestions, suggestion)
	}

	return suggestions
}

// MatchSuggestion checks if a query matches a suggestion based on the suggestion's type and a given threshold.
func MatchSuggestion(query string, suggestion TypeSuggestion, threshold int) bool {
	query = strings.ToLower(query)
	suggestionLabel := strings.ToLower(suggestion.Label)

	if query == suggestionLabel {
		return true // Short-circuit perfect match
	}
	if len(query) == 1 {
		// If query is a single character, prefer prefix match
		return strings.HasPrefix(suggestionLabel, query)
	}

	switch suggestion.Type {
	case "CreationDate":
		// For CreationDate, match if the query is numeric and is a prefix of the label
		if isNumeric(query) && strings.HasPrefix(suggestionLabel, query) {
			return true
		}
		//fuzzy matching for CreationDate if not a perfect match
		return threshold <= helpers.MatchToSuggest(query, suggestion.Label)

	case "location", "artist/band", "member":
		//fuzzy matching for locations,artist/band,members
		return threshold <= helpers.MatchToSuggest(query, suggestion.Label)

	case "firstAlbum":
		// For firstAlbum, check if it matches within a date fuzziness threshold
		if !helpers.DatesFuzz(query, suggestion.Label, threshold) {
			return threshold <= helpers.MatchToSuggest(query, suggestion.Label)
		}
		return true

	default:
		return strings.Contains(suggestionLabel, query)

	}
}

// FilterSuggestions filters a list of suggestions based on the query and matching threshold.
func FilterSuggestions(query string, suggestionsList []TypeSuggestion, threshold int) []TypeSuggestion {
	var filtered []TypeSuggestion
	for _, s := range suggestionsList {
		if MatchSuggestion(query, s, threshold) {
			filtered = append(filtered, s)
		}
	}
	return filtered
}
func isNumeric(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}
