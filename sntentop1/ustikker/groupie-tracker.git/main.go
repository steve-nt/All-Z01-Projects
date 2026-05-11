package main

import (
	"groupie-tracker/handlers"
	"groupie-tracker/middleware"
	"log"
	"net/http"
	"os"
)

func main() {
	// Configure logging
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.SetOutput(os.Stdout)

	// Create a new ServeMux to handle routes
	mux := http.NewServeMux()

	// Serve static files (CSS, JavaScript, images) from the "static" directory
	fs := http.FileServer(http.Dir("static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	// Set up routes with logging middleware
	mux.HandleFunc("/", middleware.LoggingMiddleware(handlers.LandingHandler))
	mux.HandleFunc("/main", middleware.LoggingMiddleware(handlers.IndexHandler))
	mux.HandleFunc("/geo", middleware.LoggingMiddleware(handlers.GeoHandler))

	// Start the HTTP server on localhost:8080
	log.Println("Starting server on :8080")
	if err := http.ListenAndServe("localhost:8080", mux); err != nil {
		log.Printf("Server error: %v", err)
		log.Fatalf("Server failed to start: %v", err)
	}
}
