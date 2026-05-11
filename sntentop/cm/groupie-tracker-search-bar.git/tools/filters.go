package tools

import (
	"log"
	"net/http"
	"strconv"
	"strings"
)

type tmplData struct {
	Artists          []Artist
	CreationMin      string
	CreationMax      string
	AlbumMin         string
	AlbumMax         string
	Members          []string
	Locations        []string
	SelectedLocation string
	// used string type to allow empty value. This way we don't have "0" (zero number)
	// as the default value so we don't have to delete it first and type our desired value after
}

func filters(r *http.Request) tmplData {
	var templateData tmplData
	var err error

	// 1. filter by DATE CREATED =================================================================================

	var filteredArtists1 []Artist

	creationMinStr := r.URL.Query().Get("creation-min")
	creationMaxStr := r.URL.Query().Get("creation-max")
	// this .Get always returns a string. The "number" we have used for it in html, is
	// only an input restriction for the user/client (like handling error 400)

	// convert the filter inputs to integers, handling error 400 (invalid client input) and filter the artists
	var creationMin, creationMax int

	if creationMinStr != "" {
		creationMin, err = strconv.Atoi(creationMinStr)
		if err != nil {
			creationMin = 0
		}
	}
	if creationMaxStr != "" {
		creationMax, err = strconv.Atoi(creationMaxStr)
		if err != nil {
			creationMax = 0
		}
	}

	// filtering while iterating over the artists...
	for _, artist := range artists {
		if (creationMin == 0 || artist.CreationDate >= creationMin) &&
			(creationMax == 0 || artist.CreationDate <= creationMax) {
			filteredArtists1 = append(filteredArtists1, artist)
		}
	}

	// 2. filter by FIRST ALBUM DATE =============================================================================

	var filteredArtists2 []Artist
	var albumYearStr string
	var albumDateParts []string

	albumMinStr := r.URL.Query().Get("album-min")
	albumMaxStr := r.URL.Query().Get("album-max")

	var albumMin, albumMax int

	if albumMinStr != "" {
		albumMin, err = strconv.Atoi(albumMinStr)
		if err != nil {
			albumMin = 0
		}
	}

	if albumMaxStr != "" {
		albumMax, err = strconv.Atoi(albumMaxStr)
		if err != nil {
			albumMax = 0
		}
	}

	// converting the string date of first album to int (only the year)
	// and filtering while iterating over the artists...

	// the reason we iterate over the filteredArtists2 and not artists is that even if the user didn't
	// use the filter for date created, this condition:
	// >> if (creationMin == 0 || artist.CreationDate >= creationMin) && (creationMax == 0 || artist.CreationDate <= creationMax) <<
	// will pass all the artists anyway
	for _, artist := range filteredArtists1 {

		albumDateParts = strings.Split(artist.FirstAlbum, "-")
		if len(albumDateParts) != 3 {
			log.Printf("Skipping artist %s - invalid first album date format: %s", artist.Name, artist.FirstAlbum)
			continue
		}

		albumYearStr = albumDateParts[2]

		albumYear, err := strconv.Atoi(albumYearStr)
		if err != nil {
			log.Printf("Skipping artist %s - invalid first album date format: %s", artist.Name, artist.FirstAlbum)
			continue
		}

		if (albumMin == 0 || albumYear >= albumMin) &&
			(albumMax == 0 || albumYear <= albumMax) {
			filteredArtists2 = append(filteredArtists2, artist)
		}
	}

	// 3. filter by NUMBER OF MEMBERS ============================================================================

	members := r.URL.Query()["members"]

	var filteredArtists3 []Artist

	if len(members) == 0 {
		filteredArtists3 = filteredArtists2
		// if no checkbox was selected we just use the output of filters 1 and 2 (filtfilteredArtists2)
	} else {
		for _, artist := range filteredArtists2 {
			for _, membersCheckbox := range members {

				membersCheckboxInt, err := strconv.Atoi(membersCheckbox)
				if err != nil {
					membersCheckboxInt = 0
				}

				if len(artist.Members) == membersCheckboxInt {
					filteredArtists3 = append(filteredArtists3, artist)
					break
				}
			}
		}
	}

	// 4. filter by LOCATION OF CONCERTS =========================================================================

	var filteredArtists4 []Artist

	locationMap := make(map[string]struct{})
	// we create a map to store locations w/o duplicates because in a map each key should be unique.
	// the value will be an empty struct. This way we dont have to iterate over the locations to check if one is already added.

	// with the two following iterations, we create the list of concerts locations for the drop-down menu for the
	// user to choode from

	for _, artist := range artists {
		for location := range artist.Relations {
			locationMap[location] = struct{}{}
			// with range we iterate over the keys of the artis,Relations map and store
			// each key as a new unique key in the map we created
		}
	}

	// we create a slice []string from the map keys (unique locations) for easier use in the template
	var uniqueLocations []string
	for location := range locationMap {
		uniqueLocations = append(uniqueLocations, location)
	}

	location := r.URL.Query().Get("location")

	if location != "" {
		for _, artist := range filteredArtists3 {
			for loc := range artist.Relations {
				if loc == location {
					filteredArtists4 = append(filteredArtists4, artist)
					break
				}
			}
		}
	} else {
		filteredArtists4 = filteredArtists3
	}

	// Pass the filtered artists and the filtered values to the template =========================================

	templateData.Artists = filteredArtists4
	templateData.CreationMin = creationMinStr
	templateData.CreationMax = creationMaxStr
	templateData.AlbumMin = albumMinStr
	templateData.AlbumMax = albumMaxStr
	templateData.Locations = uniqueLocations
	templateData.Members = members
	templateData.SelectedLocation = location

	return templateData

}
