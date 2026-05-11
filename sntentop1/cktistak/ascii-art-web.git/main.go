package main

import (
	"log"
	"net/http"

	"ascii-art-web/handlers"
	"ascii-art-web/services"
)

// main initializes and starts the ASCII art web server.
// It sets up the service layer, HTTP handlers, routes, and starts listening on port 8080.
func main() {
	// Initialize the ASCII art service
	asciiService := services.NewAsciiArtWeb()

	// Load banner font files from the banners directory
	if err := asciiService.LoadBanners(); err != nil {
		log.Fatalf("Failed to load banners: %v", err)
	}

	// Create HTTP handler with the initialized service
	asciiHandler := handlers.NewAsciiHandler(asciiService)

	// Register HTTP routes
	http.HandleFunc("/", asciiHandler.HandleHome)              // GET: Home page with form
	http.HandleFunc("/ascii-art", asciiHandler.HandleAsciiArt) // POST: Process ASCII art generation
	http.HandleFunc("/static/", asciiHandler.HandleResources)

	// Start the web server
	port := ":8080"
	log.Printf("Server starting on http://localhost%s", port)
	log.Println("Press Ctrl+C to stop the server")
	log.Fatal(http.ListenAndServe(port, nil))
}
