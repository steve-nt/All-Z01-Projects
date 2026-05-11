package groupie_tracker_search

import (
	"strings"
)

func Members(artists []Artist) {
	for i, artist := range artists {
		artists[i].NumMembers = len(artist.Members)

	}
}

func StrMembers(artists []Artist) {
	for i := range artists {
		artists[i].NameMembers = strings.Join(artists[i].Members, ",")

	}
}
