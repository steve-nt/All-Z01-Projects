package main

import (
	"fmt"
	"groupie-tracker-geolocalization/app/handlers"
	"groupie-tracker-geolocalization/app/services"
	"net/http"
	"os"
	"time"
)

func main() {
	if len(os.Args) != 1 {
		fmt.Fprintln(os.Stderr, "Check arguments!")
		return
	}

	// Store artists in the shared variable
	var err error
	artists, err := services.FetchArtists("https://groupietrackers.herokuapp.com/api/artists")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to fetch artists:", err)
		return
	}

	// Process artist details
	artistsWithDetails, err := services.ProcessArtistsDetails(artists)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to process artist details:", err)
		return
	}

	// Store in handlers package (through a setter function)
	handlers.SetArtistsData(artistsWithDetails)

	fmt.Println("Server running at: http://localhost:8080/")
	startTime := time.Now()
	services.LogHistory(fmt.Sprintf("Server started at %s", startTime.Format(time.RFC1123)))

	// Creates a new ServeMux instance which will be used to route HTTP requests
	mux := http.NewServeMux()
	handlers.RegisterRoutes(mux)

	// Start the HTTP server on port 8080 using the custom ServeMux
	err = http.ListenAndServe(":8080", mux)
	if err != nil {
		services.LogHistory(fmt.Sprintf("Server stopped due to error: %s", err))
		fmt.Fprintln(os.Stderr, "Server error:", err)
	}
}
