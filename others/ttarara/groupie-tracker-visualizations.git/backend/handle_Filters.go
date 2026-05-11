package backend

import (
	"bytes"
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
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Creation date range
	creationDateStart := r.FormValue("creationDateStart")
	creationDateEnd := r.FormValue("creationDateEnd")
	if creationDateStart == "" || creationDateEnd == "" {
		http.Error(w, "Missing creation date range", http.StatusBadRequest)
		return
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
	if firstAlbumDateStart == "" || firstAlbumDateEnd == "" {
		http.Error(w, "Missing first album date range", http.StatusBadRequest)
		return
	}
	albumStart, errAS := strconv.Atoi(firstAlbumDateStart)
	albumEnd, errAE := strconv.Atoi(firstAlbumDateEnd)
	if errAS != nil || errAE != nil {
		http.Error(w, "Invalid first album date range", http.StatusBadRequest)
		return
	}

	// Number-of-members (checkboxes). If user picks "5," we interpret that as "5 or more."
	numOfMembers := r.URL.Query()["Member"] // e.g. ["1","2","5"]

	// Location partial match
	locationFilter := strings.ToLower(strings.TrimSpace(r.FormValue("filter")))

	// == Fetch all artists ==
	theArtists, err := fetchData[Artists]("https://groupietrackers.herokuapp.com/api/artists")
	if err != nil {
		http.Error(w, "Failed to fetch artists data", http.StatusInternalServerError)
		return
	}

	// == Fetch location index ==
	allLocations, err := fetchData[LocationsIndex]("https://groupietrackers.herokuapp.com/api/locations")
	if err != nil {
		http.Error(w, "Failed to fetch location data", http.StatusInternalServerError)
		return
	}
	

	// == Filter logic ==
	var filtered Artists
	for _, artist := range theArtists {
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
				// If "5," treat as "5 or more"
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

		if keep {
			filtered = append(filtered, artist)
		}
	}

	// == Render the results in search_Filters.html ==
	tmpl, err := template.ParseFiles("templates/search_Filters.html")
	if err != nil {
		http.Error(w, "Failed to parse 'search_Filters.html'", http.StatusInternalServerError)
		return
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, filtered); err != nil {
		http.Error(w, "Failed to execute template", http.StatusInternalServerError)
		return
	}
	w.Write(buf.Bytes())
}


