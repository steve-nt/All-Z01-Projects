package search

import (
	"fmt"
	"groupie/fetch"
	"strconv"
	"strings"
)

type Suggestion struct {
	Display string // What to show in the dropdown
	Type    string // Entity type (member, band, location, etc.)
}

func GetSuggestions(query string, artists []fetch.Artist) []Suggestion {
	query = strings.ToLower(query)
	var suggestions []Suggestion
	seen := make(map[string]bool) // To prevent duplicates

	for _, artist := range artists {
		// Check name
		if strings.Contains(strings.ToLower(artist.Name), query) {
			key := fmt.Sprintf("name:%s", artist.Name)
			if !seen[key] {
				suggestions = append(suggestions, Suggestion{
					Display: artist.Name,
					Type:    "artist",
				})
				seen[key] = true
			}
		}

		// Check first album
		if strings.Contains(strings.ToLower(artist.FirstAlbum), query) {
			key := fmt.Sprintf("album:%s", artist.FirstAlbum)
			if !seen[key] {
				suggestions = append(suggestions, Suggestion{
					Display: artist.FirstAlbum,
					Type:    "album",
				})
				seen[key] = true
			}
		}

		// Check creation date
		if strings.Contains(strings.ToLower(strconv.Itoa(artist.CreationDate)), query) {
			key := fmt.Sprintf("date:%d", artist.CreationDate)
			if !seen[key] {
				suggestions = append(suggestions, Suggestion{
					Display: strconv.Itoa(artist.CreationDate),
					Type:    "date",
				})
				seen[key] = true
			}
		}

		// Check members
		for _, member := range artist.Members {
			if strings.Contains(strings.ToLower(member), query) {
				key := fmt.Sprintf("member:%s", member)
				if !seen[key] {
					suggestions = append(suggestions, Suggestion{
						Display: member,
						Type:    "member",
					})
					seen[key] = true
				}
			}
		}

		// Check locations
		for _, location := range artist.Locations {
			if strings.Contains(strings.ToLower(location.Name), query) {
				key := fmt.Sprintf("location:%s", location.Name)
				if !seen[key] {
					suggestions = append(suggestions, Suggestion{
						Display: location.Name,
						Type:    "location",
					})
					seen[key] = true
				}
			}
		}

	}

	// Limit to 10 suggestions
	if len(suggestions) > 10 {
		return suggestions[:10]
	}
	return suggestions
}
