package backend

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// helper function making an HTTP request, decoding JSON, handling timeouts
func fetchData[T any](apiURL string) (T, error) {
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

// Fetch artist data from given API URL
func fetchArtist(apiURL string) (Artist, error) {
	return fetchData[Artist](apiURL)
}

// Concurrently fetch multiple artists
func FetchArtists(apiURL string) (Artists, error) {
	return fetchData[Artists](apiURL)
}

// Concurrently fetch relations, locations, and concert dates
func fetchExtraDetails(artist Artist) (Artist, error) {
	relationsChan := make(chan Artist)
	locationsChan := make(chan Artist)
	datesChan := make(chan Artist)
	errorChan := make(chan error)

	// Fetch relations
	go func() {
		resp, err := http.Get(artist.Relations)
		if err != nil {
			errorChan <- err
			return
		}
		defer resp.Body.Close()

		err = json.NewDecoder(resp.Body).Decode(&artist.Relation)
		if err != nil {
			errorChan <- err
			return
		}
		relationsChan <- artist
	}()

	// Fetch locations
	go func() {
		resp, err := http.Get(artist.Locations)
		if err != nil {
			errorChan <- err
			return
		}
		defer resp.Body.Close()

		err = json.NewDecoder(resp.Body).Decode(&artist.Location)
		if err != nil {
			errorChan <- err
			return
		}
		locationsChan <- artist
	}()

	// Fetch concert dates
	go func() {
		resp, err := http.Get(artist.Dates)
		if err != nil {
			errorChan <- err
			return
		}
		defer resp.Body.Close()

		err = json.NewDecoder(resp.Body).Decode(&artist.Date)
		if err != nil {
			errorChan <- err
			return
		}
		datesChan <- artist
	}()

	for i := 0; i < 3; i++ {
		select {
		case artist = <-relationsChan:
		case artist = <-locationsChan:
		case artist = <-datesChan:
		case err := <-errorChan:
			return artist, err
		case <-time.After(5 * time.Second):
			return artist, fmt.Errorf("Timeout fetching extra details")
		}
	}
	return artist, nil
}
