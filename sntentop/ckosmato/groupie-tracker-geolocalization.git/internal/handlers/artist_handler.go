package handlers

import (
	"gtracker/internal/services"
	"log"
	"net/http"
)

// ArtistsHandler handles the request to fetch artists
func ArtistsHandler(w http.ResponseWriter, r *http.Request) {
	artists, err := services.GetArtists()
	if err != nil {
		http.Error(w, "Failed to fetch artists", http.StatusInternalServerError)
		log.Printf("Error fetching artists: %v", err)
		return
	}
	if r.URL.Path != "/" {
		// Set the status code to 404
		w.WriteHeader(http.StatusNotFound)
		// Serve a custom 404 HTML page
		if err := services.NoPageTemplate.Execute(w, artists); err != nil {
			http.Error(w, "Failed to render page", http.StatusInternalServerError)
			log.Printf("Error executing template: %v", err)
			return
		}
		return
	}

	if err := services.IndexTemplate.Execute(w, artists); err != nil {
		http.Error(w, "Failed to render page", http.StatusInternalServerError)
		log.Printf("Error executing template: %v", err)
	}
}
