package bin

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

// handleHome fetches and serves the data for the home page
func handleHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" || r.Method != http.MethodGet {
		renderErrorTemplate(w, "Page not found", generateSuggestions())
		return
	}

	filterError := r.Context().Value(errorKey)
	if filterError != nil {
		log.Printf("Error parsing filters: %v", filterError)
		renderErrorTemplate(w, "Invalid filter", generateSuggestions())
		return
	}

	artists, ok := r.Context().Value(DataKey).([]Artist)
	if !ok || len(artists) == 0 {
		log.Printf("Error fetching artists: context value is not valid or empty")
		renderErrorTemplate(w, "No artists found", generateSuggestions())
		return
	}

	totalArtists, ok := r.Context().Value(totalArtistsKey).(int)
	if !ok {
		log.Printf("Error fetching total artists count from context")
		http.Error(w, "Failed to fetch total artists count", http.StatusInternalServerError)
		return
	}

	message, ok := r.Context().Value(messageKey).(string)
	if !ok {
		log.Printf("Error fetching message from context")
		http.Error(w, "Failed to fetch message", http.StatusInternalServerError)
		return
	}

	filters, _ := parseFilters(r.URL.Query())

	data := HomePageData{
		Artists:        artists,
		Message:        message,
		Total:          totalArtists,
		Current:        len(artists),
		NextPagination: len(artists) + 12,
		Shuffle:        fmt.Sprintf("%t", r.URL.Query().Get("shuffle") == "true"),
		Filter:         filters,
		HasFilters:     len(r.URL.Query()) > 0,
		HasSearchQuery: r.URL.Query().Get("query") != "",
	}

	renderTemplate(w, "home.html", data)
}

// handleIndividualArtistRequest fetches and serves the data for an individual artist page
func handleIndividualArtistRequest(w http.ResponseWriter, r *http.Request) {
	name := strings.ReplaceAll(r.URL.Path[len("/artist/"):], "-", " ")

	cacheKey := "artist_" + name
	if data, found := artistCache.Get(cacheKey, 1*time.Minute); found {
		// Serve from cache
		renderTemplate(w, "artist.html", data)
		return
	}

	var wg sync.WaitGroup
	dataChan := make(chan ArtistPageData)
	errChan := make(chan error)

	wg.Add(1)
	go func() {
		defer wg.Done()
		artist, suggestions, err := fetchArtistByName(name)
		if err != nil {
			renderErrorTemplate(w, "Artist not found", suggestions)
			errChan <- err
			return
		}

		var relations Relation
		var dates Date
		var locations Location
		var concerts Concerts
		var maps string
		var concertsJSON string

		var fetchWg sync.WaitGroup
		fetchWg.Add(3)

		go func() {
			defer fetchWg.Done()
			var err error
			relations, err = FetchRelations(artist.Id)
			if err != nil {
				errChan <- err
				return
			}
			concerts = relations.Parse()
			concertsJSONBytes, err := json.Marshal(concerts)
			if err != nil {
				errChan <- fmt.Errorf("failed to marshal concerts: %v", err)
				return
			}
			concertsJSON = string(concertsJSONBytes)
		}()

		go func() {
			defer fetchWg.Done()
			var err error
			dates, err = FetchDates(artist.Id)
			if err != nil {
				errChan <- err
				return
			}
			dates.Parse()
		}()

		go func() {
			defer fetchWg.Done()
			var err error
			locations, err = FetchLocations(artist.Id)
			if err != nil {
				errChan <- err
				return
			}
			mapsJSONBytes, err := json.Marshal(locations)
			if err != nil {
				errChan <- fmt.Errorf("failed to marshal maps: %v", err)
				return
			}
			maps = string(mapsJSONBytes)
			locations.Parse()
		}()

		fetchWg.Wait()

		dataChan <- ArtistPageData{
			Artist:       artist,
			Dates:        dates,
			Locations:    locations,
			Relations:    concerts,
			Maps:         maps,
			ConcertsJSON: concertsJSON,
		}
	}()

	go func() {
		wg.Wait()
		close(dataChan)
		close(errChan)
	}()

	select {
	case data := <-dataChan:
		artistCache.Set(cacheKey, data, 2*time.Minute)
		renderTemplate(w, "artist.html", data)
	case err := <-errChan:
		log.Printf("Error fetching data: %v", err)
	}
}

// HandleFilterData fetches and serves the filter data to the frontend for the filter dropdowns and search suggestions
func HandleFilterData(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	cacheKey := "filter_data"
	if data, found := filterDataCache.Get(cacheKey, 2*time.Minute); found {
		w.Write(data.([]byte))
		return
	}

	artists, err := FetchAndCacheArtists()
	if err != nil {
		log.Printf("Error fetching artists: %v", err)
		http.Error(w, "Failed to fetch artists", http.StatusInternalServerError)
		return
	}

	locations, err := FetchAndCacheAllLocations()
	if err != nil {
		log.Printf("Error fetching locations: %v", err)
		http.Error(w, "Failed to fetch locations", http.StatusInternalServerError)
		return
	}
	format := TransformLocationsToFiltered(locations)

	data := struct {
		Bands     []Artist
		Locations []FilteredLocation
	}{Bands: artists, Locations: format}

	responseData, err := json.Marshal(data)
	if err != nil {
		log.Printf("Error encoding response data: %v", err)
		http.Error(w, "Failed to encode response data", http.StatusInternalServerError)
		return
	}

	filterDataCache.Set(cacheKey, responseData, 5*time.Minute)
	w.Write(responseData)
}
