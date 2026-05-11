package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"groupie-tracker/backend"
)

func main() {
	if len(os.Args) != 1 {
		fmt.Fprintln(os.Stderr, "check args!!!")
		return
	}

	// comment out to trigger 500 error and change apiArtist on PageHAndler
	backend.Init()

	fmt.Println("Server running at: http://localhost:8080/")
	startTime := time.Now()
	backend.LogHistory(fmt.Sprintf("Server started at %s", startTime.Format(time.RFC1123)))

	// Creates a new ServeMux instance which will be used to route HTTP requests
	mux := http.NewServeMux()

	// Register routes
	mux.HandleFunc("/", backend.HandleHome)
	mux.HandleFunc("/about", backend.HandleAbout)
	mux.HandleFunc("/index", backend.HandleIndex)
	mux.HandleFunc("/Artist/", backend.HandlePage)
	mux.HandleFunc("/404", backend.ErrorHandler)
	mux.Handle("/frontend/", http.StripPrefix("/frontend/", http.FileServer(http.Dir("frontend"))))
	mux.HandleFunc("/search", backend.HandleSearch)
	mux.HandleFunc("/filters", backend.HandleFilters)
	mux.HandleFunc("/all_locations", backend.HandleAllLocations)

	// Start the HTTP server on port 8080 using the custom ServeMux
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		backend.LogHistory(fmt.Sprintf("Server stopped due to error: %s", err))
		fmt.Fprintln(os.Stderr, "Server error:", err)
	}
}
