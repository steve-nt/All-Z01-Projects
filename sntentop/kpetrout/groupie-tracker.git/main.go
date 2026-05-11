package main

import (
	"fmt"
	"groupie/server"
	"net/http"
	"os"
	"time"
)

func main() {
	startNow := time.Now()
	apiURL := os.Getenv("GROUPIE_API_URL")
	if apiURL == "" {
		fmt.Println("API URL not found")
		return
	}

	// Setup routes first
	http.HandleFunc("/", server.HomeRender)
	http.HandleFunc("/homepage", server.DataRender)
	http.HandleFunc("/details", server.DetailsRender)
	http.HandleFunc("/contact", server.ContactRender)
	http.HandleFunc("/about", server.AboutRender)
	http.HandleFunc("/api/suggestions", server.SuggestionsHandler)
	http.Handle("/style/", http.StripPrefix("/style/", http.FileServer(http.Dir("style"))))
	http.Handle("/favicon.ico", http.FileServer(http.Dir("./static")))

	// Start data fetching in background
	cacheFilename := "data.json"
	server.Fetching(apiURL, cacheFilename)

	fmt.Println("Time Taken to Start:", time.Since(startNow))
	fmt.Println("Server is running on http://localhost:8080")

	// Start server
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Server error:", err)
		fmt.Println("Server is running on http://localhost:8081")

		if err := http.ListenAndServe(":8081", nil); err != nil {
			fmt.Println("Port 8081 is already in use", err)
			return
		}
	}
}
