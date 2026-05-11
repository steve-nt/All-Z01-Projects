package services

import (
	"errors"
	"gtracker/api"
	"gtracker/internal/models"
	"log"
	"strconv"
	"strings"
)

var ArtistsCache []models.Artist

func init() {

	var err error
	ArtistsCache, err = api.FetchArtists()
	if err != nil {
		log.Fatalf("Failed to fetch artists: %v", err)
	}
}

// GetArtists returns the cached artists
func GetArtists() ([]models.Artist, error) {
	if len(ArtistsCache) == 0 {
		return nil, errors.New("no artists found")
	}
	return ArtistsCache, nil
}

// FilterArtistsByCreationDate filters artists based on their creation date
func FilterArtistsByCreationDate(minYear, maxYear int) []models.Artist {
	var filteredArtists []models.Artist
	for _, artist := range ArtistsCache {
		if (minYear == 0 || artist.CreationDate >= minYear) && (maxYear == 0 || artist.CreationDate <= maxYear) {
			filteredArtists = append(filteredArtists, artist)
		}
	}
	return filteredArtists
}

// FilterArtistsByFirstAlbumDate filters artists based on the year of their first album
func FilterArtistsByFirstAlbumDate(minYear, maxYear int) []models.Artist {
	var filteredArtists []models.Artist
	for _, artist := range ArtistsCache {
		firstAlbumYear, err := extractYear(artist.FirstAlbum)
		if err != nil {
			continue
		}
		if firstAlbumYear >= minYear && firstAlbumYear <= maxYear {
			filteredArtists = append(filteredArtists, artist)
		}
	}
	return filteredArtists
}

// FilterArtistsByMembers filters artists based on the number of members
func FilterArtistsByMembers(memberCounts []int) []models.Artist {
	var filteredArtists []models.Artist
	for _, artist := range ArtistsCache {
		memberCount := len(artist.Members)
		for _, count := range memberCounts {
			if memberCount == count {
				filteredArtists = append(filteredArtists, artist)
				break
			}
		}
	}
	return filteredArtists
}

// FilterArtistsByConcertLocations filters artists based on the concert locations
func FilterArtistsByConcertLocations(locations []string) []models.Artist {
	var filteredArtists []models.Artist
	for _, artist := range ArtistsCache {
		for _, concertLocation := range artist.ConcertsFormatted {
			parts := strings.Split(concertLocation, " : ")
			if len(parts) < 2 {
				continue
			}
			location := parts[0]
			for _, loc := range locations {
				if location == loc {
					filteredArtists = append(filteredArtists, artist)
					break
				}
			}
		}
	}
	return filteredArtists
}

// Intersect returns the intersection of two artist slices
func Intersect(a, b []models.Artist) []models.Artist {
	m := make(map[int]models.Artist)
	for _, item := range a {
		m[item.Id] = item
	}
	var result []models.Artist
	for _, item := range b {
		if _, ok := m[item.Id]; ok {
			result = append(result, item)
		}
	}
	return result
}

// extractYear extracts the year from a date string in the format dd-mm-yyyy
func extractYear(dateStr string) (int, error) {
	parts := strings.Split(dateStr, "-")
	if len(parts) != 3 {
		return 0, errors.New("invalid date format")
	}
	return strconv.Atoi(parts[2])
}
