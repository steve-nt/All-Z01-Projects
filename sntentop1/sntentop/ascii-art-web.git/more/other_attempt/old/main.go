package main

import (
	"ascii-art-web/handlers"
	"html/template"
	"log"
	"net/http"
)

func main() {
	// Initialize templates and assign them to handlers.Templates.
	handlers.Templates = template.Must(template.ParseGlob("templates/*"))

	http.HandleFunc("/", handlers.HomePageHandler)
	http.HandleFunc("/ascii-art", handlers.AsciiArtHandler)

	// Serve static files (CSS and fonts).
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	log.Println("Server running at http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
