package groupie_tracker_search

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	// Set the status code to 404
	w.WriteHeader(http.StatusNotFound)
	// Path to the 404 page template
	templatePath := filepath.Join("html", "404page.html")
	// Try to parse the 404 template
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		// If there's an error loading the 404 page template, return 404 Not Found
		http.Error(w, "404 Page not found: Template could not be loaded", http.StatusNotFound)
		log.Println("Error loading 404 template:", err)
		return
	}
	// Render the template if loaded successfully
	if err := tmpl.Execute(w, nil); err != nil {
		// If there's an error rendering the template, show internal server error
		http.Error(w, "Failed to render 404 template", http.StatusInternalServerError)
		log.Println("Error rendering 404 template:", err)
	}
}
