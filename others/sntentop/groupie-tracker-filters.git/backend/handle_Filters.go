package backend

import (
	"bytes"
	"encoding/json"
	"html/template"
	"net/http"
	"strconv"
	"strings"
)

// LocationsIndex is a struct that represents the JSON structure of the locations index.
type LocationsIndex struct {
	Index []struct {
		ID        int      `json:"id"`
		Locations []string `json:"locations"`
	} `json:"index"`
}

// HandleFilters processes the user's filter inputs and renders the filtered artists
func HandleFilters(w http.ResponseWriter, r *http.Request) {
	// Only allow GET requests
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get the creation date range from the form values.
	creationDateStart := r.FormValue("creationDateStart")
	creationDateEnd := r.FormValue("creationDateEnd")

	if creationDateStart == "" {
		creationDateStart = "1958"
	}
	if creationDateEnd == "" {
		creationDateEnd = "2025"
	}
	startDate, errStart := strconv.Atoi(creationDateStart)
	endDate, errEnd := strconv.Atoi(creationDateEnd)
	if errStart != nil || errEnd != nil {
		http.Error(w, "Invalid creation date range", http.StatusBadRequest)
		return
	}

	// Get the first album date range from the form values.
	firstAlbumDateStart := r.FormValue("firstAlbumDateStart")
	firstAlbumDateEnd := r.FormValue("firstAlbumDateEnd")
	if firstAlbumDateStart == "" {
		firstAlbumDateStart = "1963"
	}
	if firstAlbumDateEnd == "" {
		firstAlbumDateEnd = "2025"
	}
	albumStart, errAS := strconv.Atoi(firstAlbumDateStart)
	albumEnd, errAE := strconv.Atoi(firstAlbumDateEnd)
	if errAS != nil || errAE != nil {
		http.Error(w, "Invalid first album date range", http.StatusBadRequest)
		return
	}

	// Get the number-of-members filter from the URL query parameters.
	numOfMembers := r.URL.Query()["Member"] // e.g. ["1","2","5"]

	// Get the location filter from the form value and normalize it.
	locationFilter := strings.ToLower(strings.TrimSpace(r.FormValue("filter")))

	// Fetch all artists from the API.
	if len(artists) == 0 {
		http.Error(w, "Artists data not loaded", http.StatusInternalServerError)
		return
	}
	theArtists := artists

	// Fetch the location index from the API.
	allLocations, err := fetchData[LocationsIndex]("https://groupietrackers.herokuapp.com/api/locations")
	if err != nil {
		http.Error(w, "Failed to fetch location data", http.StatusInternalServerError)
		return
	}

	// Filter logic: Iterate over all artists and apply the filters.
	var filtered Artists
	for _, artist := range theArtists {
		keep := true

		// Filter by creation date.
		if artist.CreationDate < startDate || artist.CreationDate > endDate {
			keep = false
		}

		// Filter by creation date.
		if keep {
			if len(artist.FirstAlbum) < 4 {
				keep = false
			} else {
				yearStr := artist.FirstAlbum[len(artist.FirstAlbum)-4:]
				yearVal, parseErr := strconv.Atoi(yearStr)
				if parseErr != nil || yearVal < albumStart || yearVal > albumEnd {
					keep = false
				}
			}
		}

		// Filter by number of members.
		if keep && len(numOfMembers) > 0 {
			matched := false
			for _, memberValue := range numOfMembers {
				n, _ := strconv.Atoi(memberValue)
				if (n == 5 && len(artist.Members) >= 5) || (len(artist.Members) == n) {
					matched = true
					break
				}
			}
			if !matched {
				keep = false
			}
		}

		// Filter by location partial match.
		if keep && locationFilter != "" {
			found := false
			// find artist's location block
			for _, locBlock := range allLocations.Index {
				if locBlock.ID == artist.ID {
					// check each location string
					for _, locStr := range locBlock.Locations {
						if strings.Contains(strings.ToLower(locStr), locationFilter) {
							found = true
							break
						}
					}
					break
				}
			}
			if !found {
				keep = false
			}
		}
		// If the artist passes all filters, add them to the filtered list.
		if keep {
			filtered = append(filtered, artist)
		}
	}

	// Render the results in the search_Filters.html template.
	tmpl, err := template.ParseFiles("templates/search_Filters.html")
	if err != nil {
		http.Error(w, "Failed to parse 'search_Filters.html'", http.StatusInternalServerError)
		return
	}
	// Execute the template with the filtered data and write the output to the response.
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, filtered); err != nil {
		http.Error(w, "Failed to execute template", http.StatusInternalServerError)
		return
	}
	w.Write(buf.Bytes())
}

// HandleMinMax επιστρέφει τα min/max των creation & album ημερομηνιών
func HandleMinMax(w http.ResponseWriter, r *http.Request) {
	minCreation, maxCreation := 9999, 0
	minAlbum, maxAlbum := 9999, 0

	for _, artist := range artists {
		if artist.CreationDate < minCreation {
			minCreation = artist.CreationDate
		}
		if artist.CreationDate > maxCreation {
			maxCreation = artist.CreationDate
		}

		if len(artist.FirstAlbum) >= 4 {
			yearStr := artist.FirstAlbum[len(artist.FirstAlbum)-4:]
			if year, err := strconv.Atoi(yearStr); err == nil {
				if year < minAlbum {
					minAlbum = year
				}
				if year > maxAlbum {
					maxAlbum = year
				}
			}
		}
	}

	response := map[string]int{
		"minCreation": minCreation,
		"maxCreation": maxCreation,
		"minAlbum":    minAlbum,
		"maxAlbum":    maxAlbum,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func HandleMaxMembers(w http.ResponseWriter, r *http.Request) {
	max := 0
	for _, artist := range artists {
		if len(artist.Members) > max {
			max = len(artist.Members)
		}
	}
	response := map[string]int{"maxMembers": max}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
