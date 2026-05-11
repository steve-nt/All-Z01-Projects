package handlers

import (
	"bytes"
	"encoding/json"
	"groupie-tracker-geolocalization/app/models"
	"groupie-tracker-geolocalization/app/services"
	"html/template"
	"net/http"
	"strconv"
	"strings"
)

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
		renderError(w, http.StatusMethodNotAllowed, r)
		return
	}

	// Creation date range
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

	// First album date range
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

	// Number-of-members (checkboxes)
	numOfMembers := r.URL.Query()["Member"]

	// Location partial match
	locationFilter := strings.ToLower(strings.TrimSpace(r.FormValue("filter")))

	// Fetch all artists
	artists := GetArtistsData()

	// Fetch location index
	allLocations, err := services.FetchData[LocationsIndex]("https://groupietrackers.herokuapp.com/api/locations")
	if err != nil {
		http.Error(w, "Failed to fetch location data", http.StatusInternalServerError)
		return
	}

	// Filter logic
	var filtered models.Artists
	for _, artist := range artists {
		keep := true

		// Creation Date
		if artist.CreationDate < startDate || artist.CreationDate > endDate {
			keep = false
		}

		// First Album Year
		if keep {
			if len(artist.FirstAlbum) < 4 {
				keep = false
			} else {
				yearStr := artist.FirstAlbum[len(artist.FirstAlbum)-4:] // e.g. "1998"
				yearVal, parseErr := strconv.Atoi(yearStr)
				if parseErr != nil || yearVal < albumStart || yearVal > albumEnd {
					keep = false
				}
			}
		}

		// Members
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

		// Location partial match
		if keep && locationFilter != "" {
			found := false
			// Find artist's location block
			for _, locBlock := range allLocations.Index {
				if locBlock.ID == artist.ID {
					// Check each location string
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

		if keep {
			filtered = append(filtered, artist)
		}
	}

	// Render the results in template
	tmpl, err := template.ParseFiles("web/templates/partials/filters.html")
	if err != nil {
		http.Error(w, "Failed to parse 'filters.html'", http.StatusInternalServerError)
		return
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, filtered); err != nil {
		http.Error(w, "Failed to execute template", http.StatusInternalServerError)
		return
	}
	w.Write(buf.Bytes())
}

// HandleMinMax returns min/max of creation & album dates
func HandleMinMax(w http.ResponseWriter, r *http.Request) {
	artists := GetArtistsData()

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
	artists := GetArtistsData()
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
