package handlers

import (
	"encoding/json"
	"fmt"
	"groupie-tracker-geolocalization/app/models"
	"groupie-tracker-geolocalization/app/services"
	"html/template"
	"net/http"
	"strconv"
	"strings"
)

// HandlePage serves the artist details page
func HandlePage(w http.ResponseWriter, r *http.Request) {
	// Extract and validate ID
	id, err := extractArtistID(r)
	if err != nil {
		renderError(w, http.StatusNotFound, r, "Invalid artist ID")
		return
	}

	// Fetch artist data
	artist, err := services.FetchArtistData(id)
	if err != nil {
		renderError(w, http.StatusInternalServerError, r, "Failed to load artist data")
		return
	}

	// Prepare template data
	data, err := prepareArtistTemplateData(artist)
	if err != nil {
		renderError(w, http.StatusInternalServerError, r, "Failed to prepare page data")
		return
	}

	// Render template
	if err := renderArtistTemplate(w, data); err != nil {
		renderError(w, http.StatusInternalServerError, r, "Failed to render page")
	}
}

// HandleIndex serves the index page with the artists
func HandleIndex(w http.ResponseWriter, r *http.Request) {
	services.LogHistory(fmt.Sprintf("Accessed Index Page - %s", r.RemoteAddr))

	// Get artists from the centralized storage
	artists := GetArtistsData()

	// Fetch and store globally if not already populated
	if len(artists) == 0 {
		renderError(w, http.StatusInternalServerError, r, "Failed to fetch artists")
		return
	}

	// Parse the index page template
	tmpl, err := template.ParseFiles("web/templates/partials/index.html")
	if err != nil {
		renderError(w, http.StatusInternalServerError, r, "Failed to load index template")
		return
	}

	// Render the index page with the artists
	if err := tmpl.Execute(w, artists); err != nil {
		renderError(w, http.StatusInternalServerError, r, "Failed to render index page")
	}
}

// serve the about page
func HandleAbout(w http.ResponseWriter, r *http.Request) {
	services.LogHistory(fmt.Sprintf("Accessed About Page - %s", r.RemoteAddr))

	if r.Method != http.MethodGet {
		renderError(w, http.StatusMethodNotAllowed, r, "Method not allowed")
		return
	}

	http.ServeFile(w, r, "web/templates/partials/about.html")
}

// Serve the home page
func HandleHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		renderError(w, http.StatusNotFound, r, "Page not found")
		return
	}
	if r.Method != http.MethodGet {
		renderError(w, http.StatusMethodNotAllowed, r, "Method not allowed")
		return
	}
	http.ServeFile(w, r, "web/templates/partials/home.html")
}

// Handle 404 error
func ErrorHandler(w http.ResponseWriter, r *http.Request) {
	renderError(w, http.StatusNotFound, r, "page not found")
}

// Extract and validate ID from URL path
func extractArtistID(r *http.Request) (int, error) {
	idStr := strings.TrimPrefix(r.URL.Path, "/Artist/")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 || id >= 53 {
		return 0, fmt.Errorf("invalid artist ID")
	}
	return id, nil
}

// Prepare data for the artist template
func prepareArtistTemplateData(artist *models.Artist) (interface{}, error) {
	// Prepare locations Json used on band HTML
	type LocationWithDate struct {
		Address string `json:"address"`
		Date    string `json:"date"`
	}

	var locs []LocationWithDate
	for loc, dates := range artist.Relation.DatesLocations {
		for _, date := range dates {
			locs = append(locs, LocationWithDate{
				Address: loc,
				Date:    date,
			})
		}
	}

	locationsJSON, err := json.Marshal(locs)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal locations with dates: %w", err)
	}

	// Prepare formatted dates locations map
	formattedDatesLocations := make(map[string][]string)
	for loc, dates := range artist.Relation.DatesLocations {
		formattedLoc := services.FormatLocation(loc)
		formattedDatesLocations[formattedLoc] = dates
	}

	return struct {
		*models.Artist
		Locations      []string
		LocationsJSON  template.JS
		DatesLocations map[string][]string
		MapScriptURL   string
	}{
		Artist:         artist,
		Locations:      artist.Location.Locations,
		LocationsJSON:  template.JS(locationsJSON),
		DatesLocations: formattedDatesLocations,
		MapScriptURL:   buildMapScriptURL(),
	}, nil

}

// Render the artist template
func renderArtistTemplate(w http.ResponseWriter, data interface{}) error {
	// Parse and execute the band details template
	tmpl, err := template.ParseFiles("web/templates/partials/band.html")
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	// Render the template with the artist data
	return tmpl.Execute(w, data)
}
