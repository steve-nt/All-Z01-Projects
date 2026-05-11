package server

import (
	"groupie/fetch"
	"strconv"
	"time"
)

type DatesAndLocations struct {
	CreationDateMin   int
	FirstAlbumDateMin int
	AllLocations      []string
}

func getDatesAndLocations(artists []fetch.Artist) DatesAndLocations {
	if len(artists) == 0 {
		return DatesAndLocations{}
	}
	year := time.Now().Year()
	creationDateMinn := year
	firstAlbumDateMinn := year
	allLocations := []string{}
	locationSet := make(map[string]bool)
	result := DatesAndLocations{}

	for _, art := range artists {
		if art.CreationDate < creationDateMinn {
			creationDateMinn = art.CreationDate
		}
		albumYear := getAlbumYear(art.FirstAlbum)
		if albumYear < firstAlbumDateMinn {
			firstAlbumDateMinn = albumYear
		}

		for _, loc := range art.Locations {
			if !locationSet[loc.Name] {
				allLocations = append(allLocations, loc.Name)
				locationSet[loc.Name] = true
			}
		}

	}

	result.CreationDateMin = creationDateMinn
	result.FirstAlbumDateMin = firstAlbumDateMinn
	result.AllLocations = allLocations
	return result
}

func getAlbumYear(album string) int {
	artistAlbumMin := album
	if len(artistAlbumMin) >= 4 {
		artistAlbumMin = artistAlbumMin[len(artistAlbumMin)-4:]
	}
	result, _ := strconv.Atoi(artistAlbumMin)
	return result
}
