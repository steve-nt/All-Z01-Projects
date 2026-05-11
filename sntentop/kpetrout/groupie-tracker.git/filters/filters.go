package filters

import (
	"groupie/fetch"
	"net/url"
	"strconv"
	"strings"
	"time"
)

func FilterArtists(artists []fetch.Artist, queryParams url.Values) []fetch.Artist {
	var filteredArtists []fetch.Artist

	year := time.Now().Year()
	creationDateMin, errMin := strconv.Atoi(queryParams.Get("creationDateMin"))
	creationDateMax, errMax := strconv.Atoi(queryParams.Get("creationDateMax"))
	if errMin != nil {
		creationDateMin = 0
	}
	if errMax != nil {
		creationDateMax = year
	}

	firstAlbumDateMin, errMinAlbum := strconv.Atoi(queryParams.Get("firstAlbumDateMin"))
	firstAlbumDateMax, errMaxAlbum := strconv.Atoi(queryParams.Get("firstAlbumDateMax"))
	if errMinAlbum != nil {
		firstAlbumDateMin = 0
	}
	if errMaxAlbum != nil {
		firstAlbumDateMax = year
	}

	locations := queryParams.Get("location")

	numberMembers := queryParams["Member"]
	var albumYear int

	for _, artist := range artists {
		if artist.CreationDate < creationDateMin || artist.CreationDate > creationDateMax {
			continue
		}

		if locations != "" {
			locationMatch := false
			for _, loc := range artist.Locations {
				if loc.Name == locations {
					locationMatch = true
					break
				}
			}
			if !locationMatch {
				continue
			}
		}

		album := strings.Split(artist.FirstAlbum, "-")
		if len(album) > 0 {
			albumYear, _ = strconv.Atoi(album[len(album)-1])
		}

		if albumYear < firstAlbumDateMin || albumYear > firstAlbumDateMax {
			continue
		}

		// Filter by number of members
		if len(numberMembers) > 0 {
			membersMatch := false
			for _, member := range numberMembers {
				memberCount, err := strconv.Atoi(member)
				if err != nil {
					continue
				}

				if memberCount == 8 {
					// Match artists with 8 or more members
					if len(artist.Members) >= 8 {
						membersMatch = true
						break
					}
				} else {
					// Exact match
					if len(artist.Members) == memberCount {
						membersMatch = true
						break
					}
				}
			}

			if !membersMatch {
				continue
			}
		}
		filteredArtists = append(filteredArtists, artist)
	}

	return filteredArtists
}
