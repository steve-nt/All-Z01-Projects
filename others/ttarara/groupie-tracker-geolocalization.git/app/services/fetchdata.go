package services

import (
	"encoding/json"
	"fmt"
	"groupie-tracker-geolocalization/app/api"
	"groupie-tracker-geolocalization/app/models"
	"net/http"
	"sync"
	"time"
)

// Declare package-level variables
var (
	locationsList []string
)

// helper function making an HTTP request, decoding JSON, handling timeouts
func FetchData[T any](apiURL string) (T, error) {
	var data T
	dataChan := make(chan T)
	errorChan := make(chan error)

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

	select {
	case data := <-dataChan:
		return data, nil
	case err := <-errorChan:
		return data, err
	case <-time.After(5 * time.Second):
		return data, fmt.Errorf("API request timed out")
	}
}

// Concurrently fetch multiple artists
func FetchArtists(apiURL string) (models.Artists, error) {
	return FetchData[models.Artists](apiURL)
}

func FetchExtraDetails(artist models.Artist) (models.Artist, error) {
	// fetch and format relations
	resp, err := http.Get(artist.Relations)
	if err != nil {
		return artist, fmt.Errorf("failed to fetch relations: %w", err)
	}
	defer resp.Body.Close()

	// Decode the JSON directly into artist.Relation
	if err := json.NewDecoder(resp.Body).Decode(&artist.Relation); err != nil {
		return artist, fmt.Errorf("failed to decode relations: %w", err)
	}

	// Format relation locations
	formattedRelations := make(map[string][]string)
	for loc, dates := range artist.Relation.DatesLocations {
		formattedRelations[FormatLocation(loc)] = dates
	}
	artist.Relation.DatesLocations = formattedRelations

	// Fetch locations
	resp, err = http.Get(artist.Locations)
	if err != nil {
		return artist, fmt.Errorf("failed to fetch locations: %w", err)
	}
	defer resp.Body.Close()

	// Decode the JSON directly into artist.Location
	if err := json.NewDecoder(resp.Body).Decode(&artist.Location); err != nil {
		return artist, fmt.Errorf("failed to decode locations: %w", err)
	}

	// Fetch dates
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

// Fetch all artists and their details on startup concurrently, preserving order.
func ArtistAndDetails(artists models.Artists) (models.Artists, error) {
	// Prepare a slice to store the updated artists with the same order
	updatedArtists := make([]models.Artist, len(artists))
	var wg sync.WaitGroup

	// Fetch extra details for all artists concurrently, preserving order
	for i, artist := range artists {
		wg.Add(1)
		go func(index int, art models.Artist) {
			defer wg.Done()
			updatedArt, err := FetchExtraDetails(art)
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

	// Update the global artists variable
	artists = updatedArtists

	// Build the location list while maintaining the order
	uniqueLocations := make(map[string]bool)
	for _, artist := range artists {
		for location := range artist.Relation.DatesLocations {
			formattedLocation := FormatLocation(location)
			if !uniqueLocations[formattedLocation] {
				uniqueLocations[formattedLocation] = true
				locationsList = append(locationsList, formattedLocation)
			}
		}
	}
	return updatedArtists, nil
}

// Use API client to fetch artist and relation data
func FetchArtistData(id int) (*models.Artist, error) {
	artistClient := api.NewArtistClient()

	// Fetch basic artist info
	artist, err := artistClient.GetArtistByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch artist: %w", err)
	}

	// Fetch relations data
	relation, err := artistClient.GetArtistRelations(id)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch relations: %w", err)
	}
	artist.Relation = *relation

	// Format locations before assigning
	formattedLocations := make([]string, 0, len(artist.Relation.DatesLocations))
	formattedRelations := make(map[string][]string)

	for loc, dates := range artist.Relation.DatesLocations {
		formattedLoc := FormatLocation(loc)
		formattedLocations = append(formattedLocations, formattedLoc)
		formattedRelations[formattedLoc] = dates
	}

	// Assign formatted data
	artist.Location = models.Locations{
		ID:        artist.ID,
		Locations: formattedLocations, // Now contains formatted locations
	}
	artist.Relation.DatesLocations = formattedRelations

	artist.Date = models.Dates{
		ID:    artist.ID,
		Dates: getAllDates(artist.Relation.DatesLocations),
	}
	return artist, nil
}

// Helper function to flatten all dates
func getAllDates(datesLocations map[string][]string) []string {
	var dates []string
	for _, dateList := range datesLocations {
		dates = append(dates, dateList...)
	}
	return dates
}

// Find the earliest date in a slice
func ParseEarliestDate(dates []string) (time.Time, error) {
	var earliest time.Time
	for i, d := range dates {
		parsed, err := time.Parse("02-01-2006", d)
		if err != nil {
			continue
		}
		if i == 0 || parsed.Before(earliest) {
			earliest = parsed
		}
	}
	return earliest, nil
}

// ProcessArtistsDetails fetches and processes details for all artists
func ProcessArtistsDetails(artists models.Artists) (models.Artists, error) {
	updatedArtists := make([]models.Artist, len(artists))
	var wg sync.WaitGroup
	var err error

	for i, artist := range artists {
		wg.Add(1)
		go func(index int, art models.Artist) {
			defer wg.Done()
			updatedArt, fetchErr := FetchExtraDetails(art)
			if fetchErr != nil {
				fmt.Printf("Failed to fetch details for artist %s: %v\n", art.Name, fetchErr)
				err = fetchErr
				updatedArtists[index] = art
				return
			}
			updatedArtists[index] = updatedArt
		}(i, artist)
	}

	wg.Wait()

	if err != nil {
		return nil, err
	}

	return updatedArtists, nil
}
