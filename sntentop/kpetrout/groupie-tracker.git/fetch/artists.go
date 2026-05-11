package fetch

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"
)

type RawArtist struct { // Struct to match the API's JSON format
	ID           int      `json:"id"`
	Image        string   `json:"image"`
	Name         string   `json:"name"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
}

type Artist struct {
	ID           int                 `json:"id"`
	Image        string              `json:"image"`
	Name         string              `json:"name"`
	Members      []string            `json:"members"`
	CreationDate int                 `json:"creationDate"`
	FirstAlbum   string              `json:"firstAlbum"`
	Locations    []*Location         `json:"locations"`
	ConcertDates []string            `json:"concertDates"`
	Relations    map[string][]string `json:"relation"`
}

func FetchArtists(URL string, completion func([]Artist)) { // Pass artists data
	artistURL := URL + "/artists"
	resp, err := http.Get(artistURL)
	if err != nil {
		fmt.Printf("Failed to fetch artists: %v", err)
		return
	}
	defer resp.Body.Close()

	var rawArtists []RawArtist
	if err := json.NewDecoder(resp.Body).Decode(&rawArtists); err != nil {
		fmt.Printf("Failed to parse artists JSON: %v", err)
		return
	}

	// Convert RawArtist to Artist
	artists := make([]Artist, len(rawArtists))
	for i, raw := range rawArtists {
		artists[i] = Artist{
			ID:           raw.ID,
			Image:        raw.Image,
			Name:         raw.Name,
			Members:      raw.Members,
			CreationDate: raw.CreationDate,
			FirstAlbum:   raw.FirstAlbum,
		}
	}

	var wg sync.WaitGroup
	updatedArtists := make([]Artist, len(artists))
	errCh := make(chan error, len(artists))

	for i, artist := range artists {

		wg.Add(1)
		go func(idx int, a Artist) {
			defer wg.Done()

			// Fetch extra details
			relationURL := URL + "/relation/" + strconv.Itoa(artist.ID)

			relations, err := fetchRelation(relationURL)
			if err != nil {
				errCh <- fmt.Errorf("error fetching relations for artist %d: %v", artist.ID, err)
				return
			}

			artist.Relations = relations
			artist.Locations = LocationParse(artist.Relations)
			artist.ConcertDates = DatesParse(artist.Relations)
			updatedArtists[i] = artist
		}(i, artist)
	}

	wg.Wait()
	close(errCh)

	var errors []error
	for err := range errCh {
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		for _, err := range errors {
			fmt.Println("Error:", err)
		}
		return
	}

	fmt.Println("Artists fetched successfully!")

	if completion != nil {
		completion(updatedArtists)
	} else {
		errCh <- fmt.Errorf("error passing the updated artist slice: %v", err)
		return
	}
}
