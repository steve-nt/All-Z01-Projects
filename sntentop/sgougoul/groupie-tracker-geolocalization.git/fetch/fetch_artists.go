package fetch

import (
	"encoding/json"
	"log"
	"net/http"
)

// Constants for external API URLs
const artistsURL = "https://groupietrackers.herokuapp.com/api/artists"

// Struct to match the artist API response
type Artist struct {
	ID           int      `json:"id"`
	Image        string   `json:"image"`
	Name         string   `json:"name"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
	Locations    string   `json:"locations"`    // URL for locations
	ConcertDates string   `json:"concertDates"` // URL for concert dates
	Relations    string   `json:"relations"`    // URL for relations
}

// Fetch artists from the external API
func FetchArtists() ([]Artist, error) {
	resp, err := http.Get(artistsURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var artists []Artist
	if err := json.NewDecoder(resp.Body).Decode(&artists); err != nil {
		log.Printf("%v", err)
		return nil, err
	}
	// Log the first album year for each artist
	for _, artist := range artists {
		// Extract the year from the first album (assuming it's the first 4 characters of the string)
		albumYear := ""
		if len(artist.FirstAlbum) >= 4 {
			albumYear = artist.FirstAlbum[6:10] // Get the first 4 characters (year)
		}

		// Log the artist's name and the extracted year
		log.Printf("Artist: %s, First Album Year: %s", artist.Name, albumYear)
	}

	return artists, nil
}
