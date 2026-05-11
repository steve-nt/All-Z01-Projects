package api

import (
	"fmt"
	"groupie/utils"
	"log"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"text/template"
)

// handler for main page
func MainHandler(writer http.ResponseWriter, request *http.Request) {
	if request.URL.Path != "/" {
		SendErrorPage(writer, 404, "404 - Page not found")
		return
	}

	writer.Header().Set("Cache-Control", "public, max-age=3600")

	err := request.ParseForm()
	if err != nil {
		SendErrorPage(writer, 400, "400 - Bad request")
		return
	}

	//setting default values:
	filter := FilterT{
		CreationYearStart:   1956,
		CreationYearEnd:     2025,
		FirstAlbumYearStart: 1956,
		FirstAlbumYearEnd:   2025,
		ConcertFilter:       "any"}

	tempBandSizeSlice := request.Form["band_size"]
	if len(tempBandSizeSlice) == 0 {
		for i := range filter.BandSizeFilterCheckboxes {
			filter.BandSizeFilterCheckboxes[i] = 1
			filter.BandSizeFilter = append(filter.BandSizeFilter, 1+i)
		}
	} else {
		for _, checkboxIndexStr := range tempBandSizeSlice {
			checkboxIndex, err := strconv.Atoi(checkboxIndexStr)
			if err != nil || checkboxIndex <= 0 || checkboxIndex > 10 {
				SendErrorPage(writer, 400, "400 - Bad Request <br><br> Bad values in the URL")
				return
			}
			filter.BandSizeFilterCheckboxes[checkboxIndex-1] = 1 // 1 means that it's on
			filter.BandSizeFilter = append(filter.BandSizeFilter, checkboxIndex)
		}
	}

	temp := request.FormValue("creation_year_start")
	if temp != "" {
		if filter.CreationYearStart, err = strconv.Atoi(temp); err != nil {
			SendErrorPage(writer, 400, "400 - Bad Request <br><br> Bad values in the URL")
			return
		}
	}

	temp = request.FormValue("creation_year_end")
	if temp != "" {
		if filter.CreationYearEnd, err = strconv.Atoi(temp); err != nil {
			SendErrorPage(writer, 400, "400 - Bad Request <br><br> Bad values in the URL")
			return
		}
	}
	temp = request.FormValue("first_album_year_start")
	if temp != "" {
		if filter.FirstAlbumYearStart, err = strconv.Atoi(temp); err != nil {
			SendErrorPage(writer, 400, "400 - Bad Request <br><br> Bad values in the URL")
			return
		}
	}
	temp = request.FormValue("first_album_year_end")
	if temp != "" {
		if filter.FirstAlbumYearEnd, err = strconv.Atoi(temp); err != nil {
			SendErrorPage(writer, 400, "400 - Bad Request <br><br> Bad values in the URL")
			return
		}
	}

	temp = request.FormValue("concert-filter")
	if temp != "" {
		filter.ConcertFilter = temp
	}

	filter.SearchBar = request.FormValue("searchbar")
	SearchBar := filter.SearchBar

	// Parse the artistID from the query parameters
	artistIDStr := request.URL.Query().Get("artistID")

	artistID := 0 //default value, 0 should display nothing on the right side

	if artistIDStr != "" {
		artistID, err = strconv.Atoi(artistIDStr)
		if err != nil || artistID < -1 || artistID > len(Artists) {
			SendErrorPage(writer, 400, "400 - Bad Request <br><br> Bad values in the URL")
			return
		}
	}

	// Find artist by ID
	selectedArtist, selectedArtistFound := ArtistMap[artistID]
	selectedRelation, selectedRelationFound := ArtistRelationMap[artistID]

	// If the artist is not found, return a 404 error
	if !selectedArtistFound && artistID != 0 {
		SendErrorPage(writer, 404, "404 - Artist not found")
		return
	}

	// If relation data is not found, return a 404 error
	if !selectedRelationFound && artistID != 0 {
		SendErrorPage(writer, 404, "404 - Relation data not found")
		return
	}

	selectedRelation.DatesLocations = utils.FormatMapKeys(selectedRelation.DatesLocations)
	selectedRelation.DatesLocations = utils.SortDates(selectedRelation.DatesLocations)
	sortedLocations := utils.SortLocations(selectedRelation.DatesLocations)

	//get all locations to send them for location dropdown filter
	allLocationsMap := make(map[string][]string)
	for _, relation := range ArtistRelationMap {
		for location := range utils.FormatMapKeys(relation.DatesLocations) {
			allLocationsMap[location] = []string{}
		}
	}

	//default options to show all locations
	allLocations := []string{"any"}
	allLocations = append(allLocations, utils.SortLocations(allLocationsMap)...)

	// Filter artists by filters
	filterReducedArtists := filterArtists(filter, Artists)

	// Filter artists by search query
	searchReducedArtists := searchFilter(SearchBar, filterReducedArtists)

	data := struct {
		Artists          []Artist
		Artist           Artist
		Relation         Relation
		SelectedArtistID int
		SortedLocations  []string
		AllLocations     []string
		Filter           FilterT
	}{
		Artists:          searchReducedArtists,
		Artist:           selectedArtist,
		Relation:         selectedRelation,
		SelectedArtistID: artistID,
		SortedLocations:  sortedLocations,
		AllLocations:     allLocations,
		Filter:           filter,
	}

	tmpl, err := template.ParseFiles("../templates/index.html")
	if err != nil {
		SendErrorPage(writer, 500, "500 - Internal Server Error")
		return
	}

	if err := tmpl.Execute(writer, data); err != nil {
		log.Printf("Template execution error: %v", err)
	}
}

// returns artists that match all through all the filters
func filterArtists(filter FilterT, artists []Artist) []Artist {
	newArtistSlice := []Artist{}
	for _, artist := range artists {

		//BAND SIZE FILTER
		if !slices.Contains(filter.BandSizeFilter, len(artist.Members)) {
			continue
		}

		//CREATION YEAR FILTER
		if artist.CreationDate > filter.CreationYearEnd || artist.CreationDate < filter.CreationYearStart {
			continue
		}

		//FIRST ALBUM FILTER
		fullDate := fmt.Sprint(artist.FirstAlbum)
		parts := strings.Split(fullDate, "-")
		year := parts[2]
		if year > fmt.Sprint(filter.FirstAlbumYearEnd) || year < fmt.Sprint(filter.FirstAlbumYearStart) {
			continue
		}

		//CONCERT FILTER
		if filter.ConcertFilter != "any" {

			found := false
			for key := range ArtistRelationMap[artist.ID].DatesLocations {
				if utils.FixKey(key) == filter.ConcertFilter {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}

		newArtistSlice = append(newArtistSlice, artist)
	}

	return newArtistSlice
}

// returns artists that match at least one of the search criteria
func searchFilter(searchText string, artists []Artist) []Artist {
	if searchText == "" {
		return artists
	}
	newArtistSlice := []Artist{}
	for _, artist := range artists {
		_, ok := searchMatch(searchText, artist)
		if ok {
			newArtistSlice = append(newArtistSlice, artist)
		}
	}
	return newArtistSlice
}

// returns which artist fields match the search query
func searchMatch(searchWord string, artist Artist) ([]string, bool) {

	doFlags := struct {
		artist     bool
		member     bool
		creation   bool
		firstAlbum bool
		concert    bool
	}{}

	//trim query in case there's stuff left from suggestion autocomplete
	searchWord = strings.TrimSpace(searchWord)

	//check for suffixes, and set up flags to only find results that these tags demand
	switch {
	case strings.HasSuffix(searchWord, " - artist/band"):
		doFlags.artist = true
		searchWord = strings.TrimSuffix(searchWord, " - artist/band")

	case strings.HasSuffix(searchWord, " - member"):
		doFlags.member = true
		searchWord = strings.TrimSuffix(searchWord, " - member")

	case strings.HasSuffix(searchWord, " - creation date"):
		doFlags.creation = true
		searchWord = strings.TrimSuffix(searchWord, " - creation date")

	case strings.HasSuffix(searchWord, " - first album"):
		doFlags.firstAlbum = true
		searchWord = strings.TrimSuffix(searchWord, " - first album")

	case strings.HasSuffix(searchWord, " - concert location"):
		doFlags.concert = true
		searchWord = strings.TrimSuffix(searchWord, " - concert location")

	default:
		doFlags.artist = true
		doFlags.member = true
		doFlags.creation = true
		doFlags.firstAlbum = true
		doFlags.concert = true
	}

	searchWord = strings.TrimSpace(searchWord)

	//all the matches that matched this particular artist
	matches := []string{}

	// check full artist name

	if doFlags.artist && isArtistMatch(searchWord, artist) {
		matches = append(matches, fmt.Sprint(artist.Name+" - artist/band"))
	}

	//member
	if doFlags.member {
		matchesMembers, ok := areMembersMatch(searchWord, artist)
		if ok {
			matches = append(matches, matchesMembers...)
		}
	}

	//creation date
	if doFlags.creation && isCreationDateMatch(searchWord, artist) {
		matches = append(matches, fmt.Sprint(artist.CreationDate, " - creation date"))
	}

	//first album
	if doFlags.firstAlbum && isFirstAlbumMatch(searchWord, artist) {
		matches = append(matches, fmt.Sprint(artist.FirstAlbum, " - first album"))
	}

	//check first album date
	if doFlags.concert {
		someMatches, ok := isLocationMatch(searchWord, artist)
		if ok {
			matches = append(matches, someMatches...)
		}
	}

	match := len(matches) > 0
	return matches, match
}

func isArtistMatch(searchWord string, artist Artist) bool {
	if utils.SameEnough(artist.Name, searchWord) {
		//MATCH FOUND
		return true

	} else {
		//check each "word" of artist name
		for _, namePart := range utils.SplitByWords(artist.Name) {
			if utils.SameEnough(namePart, searchWord) {
				//MATCH FOUND
				return true
			}
		}
	}
	return false
}

func areMembersMatch(searchWord string, artist Artist) ([]string, bool) {

	matches := []string{}
nextMember:
	for _, memberName := range artist.Members {
		//check full member name
		if utils.SameEnough(memberName, searchWord) {
			//MATCH FOUND
			matches = append(matches, fmt.Sprint(memberName+" - member"))
			continue

		} else {
			//check each word of member name
			for _, namePart := range utils.SplitByWords(memberName) {
				if utils.SameEnough(namePart, searchWord) {
					//MATCH FOUND
					matches = append(matches, fmt.Sprint(memberName+" - member"))
					continue nextMember
				}
			}
		}
	}
	return matches, len(matches) > 0
}
func isFirstAlbumMatch(searchWord string, artist Artist) bool {
	if utils.SameEnough(artist.FirstAlbum, searchWord) {
		//MATCH FOUND
		return true
	} else {
		//check each part of the first album date, day, month, year
		for _, datePart := range utils.SplitByWords(artist.FirstAlbum) {
			if utils.SameEnough(datePart, searchWord) {
				//MATCH FOUND
				return true
			}
		}
	}
	return false
}

func isCreationDateMatch(searchWord string, artist Artist) bool {
	if utils.SameEnough(fmt.Sprint(artist.CreationDate), searchWord) {
		//MATCH FOUND
		return true
	}
	return false
}

func isLocationMatch(searchWord string, artist Artist) ([]string, bool) {
	matches := []string{}
	//check each concert location
nextLocation:
	for location := range ArtistRelationMap[artist.ID].DatesLocations {

		//check full concert location name
		if utils.SameEnough(location, searchWord) {
			//MATCH FOUND
			matches = append(matches, fmt.Sprint(location, " - concert location"))
			continue

			//check full concert location name bug with location format same as what is visually displayed to user
		} else if utils.SameEnough(utils.FixKey(location), searchWord) {
			matches = append(matches, fmt.Sprint(utils.FixKey(location), " - concert location"))
			continue

			//checl full concert name, unformatted, but replace "-" with " ", so instad of "osaka-japan" it's "osaka japan"
		} else if utils.SameEnough(strings.Replace(location, "-", " ", -1), searchWord) {
			matches = append(matches, fmt.Sprint(strings.Replace(location, "-", " ", -1), " - concert location"))
			continue

		} else {
			//check each part of concert name
			for _, locationPart := range utils.SplitByWords(location) {
				if utils.SameEnough(locationPart, searchWord) {
					//MATCH FOUND
					matches = append(matches, fmt.Sprint(location, " - concert location"))
					continue nextLocation
				}
			}

			//check each part of concert name, formatted the same way as is displayed to user
			for _, locationPart := range utils.SplitByWords(utils.FixKey(location)) {
				if utils.SameEnough(locationPart, searchWord) {
					//MATCH FOUND
					matches = append(matches, fmt.Sprint(utils.FixKey(location), " - concert location"))
					continue nextLocation
				}
			}

			//check each part of concert name, with "-" replaced with " "
			for _, locationPart := range utils.SplitByWords(strings.Replace(location, "-", " ", -1)) {
				if utils.SameEnough(locationPart, searchWord) {
					//MATCH FOUND
					matches = append(matches, fmt.Sprint(strings.Replace(location, "-", " ", -1), " - concert location"))
					continue nextLocation
				}
			}

		}
	}

	return matches, len(matches) > 0
}
