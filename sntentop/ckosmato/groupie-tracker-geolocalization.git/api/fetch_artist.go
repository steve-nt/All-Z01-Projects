package api

import (
	"encoding/json"
	"gtracker/internal/models"
	"gtracker/utils"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	ApiBaseUrl string
	client     = &http.Client{Timeout: 5 * time.Second}
)

// FetchArtists fetches the artists from the groupietrackers API
// and concurrently fetches the concerts for each artist
func FetchArtists() ([]models.Artist, error) {
	url := "https://groupietrackers.herokuapp.com/api/artists"
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	res, getErr := client.Do(req)
	if getErr != nil {
		return nil, getErr
	}
	if res.Body != nil {
		defer res.Body.Close()
	}

	body, readErr := io.ReadAll(res.Body)
	if readErr != nil {
		return nil, readErr
	}
	var artists []models.Artist
	if err := json.Unmarshal(body, &artists); err != nil {
		return nil, err
	}

	const workerCount = 5

	var wg sync.WaitGroup

	// Create a channel to send artists
	jobs := make(chan *models.Artist, len(artists))

	for i := 0; i < workerCount; i++ {
		go func() {
			for artist := range jobs {
				FetchConcerts(artist, client, &wg)
			}
		}()
	}

	for i := range artists {
		wg.Add(1)
		jobs <- &artists[i]
	}

	// Close the jobs channel after all artists have been dispatched for concert fetching
	close(jobs)
	wg.Wait()
	// Wait for all workers to finish processing before transforming the data
	for i := range artists {
		// Perform transformation on each artist after their concerts have been fetched
		TransformArtist(&artists[i])
	}
	return artists, nil
}

// FetchConcerts fetches the concerts for a given artist
func FetchConcerts(artist *models.Artist, client *http.Client, wg *sync.WaitGroup) {
	defer wg.Done()
	artistID := strconv.Itoa(artist.Id)
	concertUrl := "https://groupietrackers.herokuapp.com/api/relation/" + artistID

	req, err := http.NewRequest(http.MethodGet, concertUrl, nil)
	if err != nil {
		log.Fatal(err)
	}
	res, getErr := client.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}
	if res.Body != nil {
		defer res.Body.Close()
	}

	body, readErr := io.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	var tempData struct {
		Concerts map[string][]string `json:"datesLocations"`
	}

	if err := json.Unmarshal(body, &tempData); err != nil {
		log.Fatal(err)
	}

	artist.Concerts = tempData.Concerts

}

// TransformArtist transforms the artist data to a more readable format
func TransformArtist(Artist *models.Artist) {

	for _, member := range Artist.Members {
		Artist.MembersFormatted = append(Artist.MembersFormatted, member+",")
	}

	for key, value := range Artist.Concerts {
		key = utils.Replacer(key)
		keySl := strings.Split(key, " ")

		for i, word := range keySl {
			if len(word) > 0 {
				if len(word) < 4 && (strings.HasPrefix(word, "uk") || strings.HasPrefix(word, "usa")) {
					keySl[i] = strings.ToUpper(word)
				} else {
					keySl[i] = strings.ToUpper(word[:1]) + word[1:]
				}
			}
		}
		key = strings.Join(keySl, " ")
		Artist.ConcertsFormatted = append(Artist.ConcertsFormatted, key+" : "+strings.Join(value, ", ")+"\n")
	}

}
