package groupie_tracker_search

import (
	"html/template"
	"log"
	"net/http"
	"strconv"
)

func ArtistPage(w http.ResponseWriter, r *http.Request) {
	// Get the artist ID from the query parameter (e.g., ?id=1)
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("Invalid artist ID provided: %s", idStr)
		http.Error(w, "Invalid artist ID", http.StatusBadRequest)
		return
	}

	// Fetch all artists using FetchArtists() function
	artists, err := getCachedArtists()
	if err != nil {
		log.Printf("Error fetching artists: %v", err)
		http.Error(w, "Failed to fetch artists", http.StatusInternalServerError)
		return
	}

	// Find the artist with the matching ID
	var selectedArtist *Artist
	for i := range artists {
		if artists[i].ID == id {
			selectedArtist = &artists[i]

			break
		}
	}

	// If no artist found, return a 404 error
	if selectedArtist == nil {
		log.Printf("Artist not found for ID: %d", id)
		NotFoundHandler(w, r) // Serve the custom 404 HTML page
		return
	}

	// Parse the template to render artist details
	tmplPath := "html/href.html" // Adjust the path as needed
	//log.Printf("Trying to load template from: %s", tmplPath)
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		log.Printf("Error loading template: %v", err)
		http.Error(w, "Unable to load template", http.StatusInternalServerError)
		return
	}

	// Render the template and pass the selected artist data
	err = tmpl.Execute(w, selectedArtist)
	if err != nil {
		log.Printf("Error rendering template for artist ID %d: %v", id, err)
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
	}
}
