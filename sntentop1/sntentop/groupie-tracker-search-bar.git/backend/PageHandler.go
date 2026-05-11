package backend

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
)

func HandlePage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		http.ServeFile(w, r, "templates/405.html")
		return
	}

	id := strings.TrimPrefix(r.URL.Path, "/Artist/")
	idd, err := strconv.Atoi(id)
	if err != nil || idd <= 0 || idd >= 53 {
		http.Redirect(w, r, "/404", http.StatusFound)
		return
	}

	apiURL := "https://groupietrackers.herokuapp.com/api/artists/" + id
	artist, err := fetchData[Artist](apiURL)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	artist, err = fetchExtraDetails(artist)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	//  Sorting locations by dates
	type LocationWithDate struct {
		Location  string
		FirstDate string
		Parsed    time.Time
	}

	var sorted []LocationWithDate
	for loc, dates := range artist.Relation.DatesLocations {
		if len(dates) == 0 {
			continue
		}
		parsed, err := time.Parse("02-01-2006", dates[0])
		if err != nil {
			continue
		}
		sorted = append(sorted, LocationWithDate{
			Location:  loc,
			FirstDate: dates[0],
			Parsed:    parsed,
		})
	}

	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Parsed.Before(sorted[j].Parsed)
	})

	var locations []string
	for _, item := range sorted {
		locations = append(locations, item.Location)
	}

	mapScriptURL := fmt.Sprintf(
		"https://maps.googleapis.com/maps/api/js?key=%s&libraries=places,marker&loading=async&map_ids=500755a5e04d8f95&callback=initMap",
		"AIzaSyDtuaYfzbNShwjWrDwBkEnhp2H3Jq9aG9g", // API key (could be in environment variable for security)
	)

	locationsJSON, err := json.Marshal(locations)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to encode locations"))
		return
	}

	data := struct {
		Artist
		Locations     []string
		LocationsJSON template.JS
		MapScriptURL  string
	}{
		Artist:        artist,
		Locations:     locations,
		LocationsJSON: template.JS(locationsJSON),
		MapScriptURL:  mapScriptURL,
	}

	tmpl, err := template.ParseFiles("templates/band.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		http.ServeFile(w, r, "templates/500.html")
		return
	}

	tmpl.Execute(w, data)
}

func earliestDate(dates []string) (time.Time, error) {
	var earliest time.Time
	for i, d := range dates {
		parsed, err := time.Parse("02-01-2006", d)
		if err != nil {
			continue
		}
		if i == 0 || parsed.Before(earliest) {
			earliest = parsed
		}
	}
	return earliest, nil
}

// HandlePage serves the artist details page.
func HandleIndex(w http.ResponseWriter, r *http.Request) {
	LogHistory(fmt.Sprintf("Accessed Index Page - %s", r.RemoteAddr))

	// Trigger error 500 by changing the link and commenting the func Init in main.go
	apiArtist := "https://groupietrackers.herokom/api/artists"

	// Fetch and store globally if not already populated
	if len(artists) == 0 {
		fetchedArtists, err := FetchArtists(apiArtist)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			http.ServeFile(w, r, "templates/500.html")
			return
		}
		artists = fetchedArtists // Assign to global variable
	}

	// Parse the index page template
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		http.ServeFile(w, r, "templates/500.html")
		return
	}

	// Render the index page with the artists
	tmpl.Execute(w, artists)
}

// serve the about page
func HandleAbout(w http.ResponseWriter, r *http.Request) {
	LogHistory(fmt.Sprintf("Accessed About Page - %s", r.RemoteAddr))
	http.ServeFile(w, r, "templates/about.html")
}
