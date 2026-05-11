// main.go
package main

import (
	"ascii-art-web/handlers"
	"log"
	"net/http"
)

func main() {
	// Serve static files (CSS, JS, images)
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Route handling
	http.HandleFunc("/", handlers.HomePageHandler)
	http.HandleFunc("/ascii-art", handlers.AsciiArtHandler)

	// Start the server
	log.Fatal(http.ListenAndServe(":8081", nil))
}
