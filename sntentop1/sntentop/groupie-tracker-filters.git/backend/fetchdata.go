package backend

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

// fetchData is a generic helper function that makes an HTTP GET request to the provided API URL,
// decodes the JSON response into the specified type T, and handles timeouts.
func fetchData[T any](apiURL string) (T, error) {
	var data T
	dataChan := make(chan T)
	errorChan := make(chan error)

	// Launch a goroutine to perform the HTTP request and JSON decoding.
	go func() {
		resp, err := http.Get(apiURL)
		if err != nil {
			errorChan <- err
			return
		}
		defer resp.Body.Close()

		err = json.NewDecoder(resp.Body).Decode(&data)
		if err != nil {
			errorChan <- err
			return
		}
		dataChan <- data
	}()
	// Use a select statement to handle the channels and timeouts.
	select {
	case data := <-dataChan:
		return data, nil
	case err := <-errorChan:
		return data, err
	case <-time.After(5 * time.Second):
		return data, fmt.Errorf("API request timed out")
	}
}

// FetchArtists is a function that fetches a list of artists from the provided API URL.
func FetchArtists(apiURL string) (Artists, error) {
	return fetchData[Artists](apiURL)
}

// fetchExtraDetails is a function that fetches additional details for a given artist.
func fetchExtraDetails(artist Artist) (Artist, error) {
	resp, err := http.Get(artist.Relations)
	if err != nil {
		return artist, fmt.Errorf("failed to fetch relations: %w", err)
	}
	defer resp.Body.Close()

	// Decode the JSON directly into artist.Relation
	if err := json.NewDecoder(resp.Body).Decode(&artist.Relation); err != nil {
		return artist, fmt.Errorf("failed to decode relations: %w", err)
	}

	// Format the keys in the DatesLocations map
	formattedDatesLocations := make(map[string][]string)
	for location, dates := range artist.Relation.DatesLocations {
		formattedLocation := FormatLocation(location)
		formattedDatesLocations[formattedLocation] = dates
	}
	artist.Relation.DatesLocations = formattedDatesLocations

	resp, err = http.Get(artist.Locations)
	if err != nil {
		return artist, fmt.Errorf("failed to fetch locations: %w", err)
	}
	defer resp.Body.Close()

	// Decode the JSON directly into artist.Location
	if err := json.NewDecoder(resp.Body).Decode(&artist.Location); err != nil {
		return artist, fmt.Errorf("failed to decode locations: %w", err)
	}

	// Format the locations using the FormatLocation function from utilities.go
	for i, location := range artist.Location.Locations {
		artist.Location.Locations[i] = FormatLocation(location)
	}

	resp, err = http.Get(artist.Dates)
	if err != nil {
		return artist, fmt.Errorf("failed to fetch dates: %w", err)
	}
	defer resp.Body.Close()

	// Decode the JSON directly into artist.Date
	if err := json.NewDecoder(resp.Body).Decode(&artist.Date); err != nil {
		return artist, fmt.Errorf("failed to decode dates: %w", err)
	}

	return artist, nil
}

// Init fetches all artists and their details on startup concurrently, preserving order.
func Init() {
	apiArtist := "https://groupietrackers.herokuapp.com/api/artists"
	var err error

	// Fetch the initial list of artists
	artists, err = FetchArtists(apiArtist)
	if err != nil {
		fmt.Println("ERROR: Failed to load artists on startup:", err)
		return
	}

	// Prepare a slice to store the updated artists with the same order
	updatedArtists := make([]Artist, len(artists))
	var wg sync.WaitGroup

	// Fetch extra details for all artists concurrently, preserving order
	for i, artist := range artists {
		wg.Add(1)
		go func(index int, art Artist) {
			defer wg.Done()
			updatedArt, err := fetchExtraDetails(art)
			if err != nil {
				fmt.Printf("Failed to fetch details for artist %s: %v\n", art.Name, err)
				updatedArtists[index] = art // fallback to original if fetch fails
				return
			}
			updatedArtists[index] = updatedArt
		}(i, artist)
	}

	// Wait for all goroutines to finish
	wg.Wait()

	// Build the location list while maintaining the order
	uniqueLocations := make(map[string]bool)
	for _, art := range updatedArtists {
		for location := range art.Relation.DatesLocations {
			formattedLocation := FormatLocation(location)
			if !uniqueLocations[formattedLocation] {
				uniqueLocations[formattedLocation] = true
				locationsList = append(locationsList, formattedLocation)
			}
		}
	}

	// Assign the correctly ordered artists
	artists = updatedArtists
}
