package filters

import (
	"strconv"
	"strings"
	"time"

	"sgougoupractice/fetch"
)

// Struct για τα φίλτρα από το frontend
type FilterOptions struct {
	CreationDateRange   [2]int
	FirstAlbumYearRange [2]int
	MemberCounts        []int
	Locations           []string
}

// Λειτουργία για φιλτράρισμα όλων των artists
func FilterArtists(artists []fetch.Artist, opts FilterOptions) []fetch.Artist {
	var filtered []fetch.Artist

	for _, artist := range artists {
		if !passesCreationDate(artist, opts.CreationDateRange) {
			continue
		}
		if !passesFirstAlbumYear(artist, opts.FirstAlbumYearRange) {
			continue
		}
		if !passesMemberCount(artist, opts.MemberCounts) {
			continue
		}
		if !passesLocation(artist, opts.Locations) {
			continue
		}
		filtered = append(filtered, artist)
	}
	return filtered
}

// --- FILTER HELPERS --- //

func passesCreationDate(artist fetch.Artist, r [2]int) bool {
	return artist.CreationDate >= r[0] && artist.CreationDate <= r[1]
}

func passesFirstAlbumYear(artist fetch.Artist, r [2]int) bool {
	year := extractYear(artist.FirstAlbum)
	return year >= r[0] && year <= r[1]
}

// Extracts year from "dd-mm-yyyy" format or from user-provided year (yyyy).
func extractYear(date string) int {
	// Check if the date is just a year (yyyy)
	if len(date) == 4 {
		year, err := strconv.Atoi(date)
		if err == nil {
			return year
		}
		return 0 // Invalid year format
	}

	// Otherwise, parse the full date in "dd-mm-yyyy" format
	t, err := time.Parse("02-01-2006", date)
	if err != nil {
		// Fallback: manually extract year from "dd-mm-yyyy" format
		parts := strings.Split(date, "-")
		if len(parts) == 3 {
			y, err := strconv.Atoi(parts[2])
			if err == nil {
				return y
			}
		}
		return 0
	}

	return t.Year() // Return the year from the parsed date
}

func passesMemberCount(artist fetch.Artist, allowed []int) bool {
	count := len(artist.Members)
	for _, v := range allowed {
		if count == v || (v >= 7 && count >= 4) {
			return true
		}
	}
	return len(allowed) == 0
}

func passesLocation(artist fetch.Artist, allowed []string) bool {
	if len(allowed) == 0 {
		return true
	}
	for _, location := range allowed {
		if strings.Contains(strings.ToLower(artist.Locations), strings.ToLower(location)) {
			return true
		}
	}
	return false
}
