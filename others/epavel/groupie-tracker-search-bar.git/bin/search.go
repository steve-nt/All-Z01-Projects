package bin

import (
	"fmt"
	"strings"
)

// SearchArtists searches for artists that match the query
func SearchArtists(artists []Artist, query string) []Artist {
	var filteredArtists []Artist
	for _, artist := range artists {
		if strings.Contains(strings.ToLower(artist.Name), query) ||
			strings.Contains(strings.ToLower(artist.FirstAlbum), query) ||
			strings.Contains(strings.ToLower(fmt.Sprintf("%d", artist.StartYear)), query) ||
			strings.Contains(artist.FirstAlbum, query) {
			filteredArtists = append(filteredArtists, artist)
			continue
		}
		for _, member := range artist.Members {
			if strings.Contains(strings.ToLower(member), query) {
				filteredArtists = append(filteredArtists, artist)
				break
			}
		}
		for _, location := range AllLocations[artist.Id-1].Locations {
			if strings.Contains(strings.ToLower(location), query) {
				filteredArtists = append(filteredArtists, artist)
				break
			}
		}
	}
	return filteredArtists
}
