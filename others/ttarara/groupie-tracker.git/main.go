package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"groupie-tracker/backend"
)

var logChan = make(chan string, 100)

func init() {
	go func() {
		logFile, err := os.OpenFile("history.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Failed to open log file:", err)
			return
		}
		defer logFile.Close()

		logger := log.New(logFile, "", log.LstdFlags)
		for msg := range logChan {
			logger.Println(msg)
		}
	}()
}

func logHistory(message string) {
	logChan <- message
}

func main() {
	if len(os.Args) != 1 {
		fmt.Fprintln(os.Stderr, "check args!!!")
		return
	}

	fmt.Println("Server running at: http://localhost:8080/")

	startTime := time.Now()
	logHistory(fmt.Sprintf("Server started at %s", startTime.Format(time.RFC1123)))

	// Serve home.html at "/"
	http.HandleFunc("/", backend.HandleHome)

	// Serve index.html at "/index" and load artist data
	http.HandleFunc("/index", func(w http.ResponseWriter, r *http.Request) {
		logHistory(fmt.Sprintf("Accessed Index Page - %s", r.RemoteAddr))

		apiArtist := "https://groupietrackers.herokuapp.com/api/artists"
		artists, err := backend.FetchArtists(apiArtist)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			http.ServeFile(w, r, "templates/500.html")
			return
		}

		tmpl, err := template.ParseFiles("templates/index.html")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			http.ServeFile(w, r, "templates/500.html")
			return
		}

		tmpl.Execute(w, artists) // Send data to the page
	})

	// Serve about.html at "/about"
	http.HandleFunc("/about", func(w http.ResponseWriter, r *http.Request) {
		logHistory(fmt.Sprintf("Accessed About Page - %s", r.RemoteAddr))
		http.ServeFile(w, r, "templates/about.html")
	})

	// Serve artist pages
	http.HandleFunc("/Artist/", backend.HandlePage)

	// Serve error pages
	http.HandleFunc("/404", backend.ErrorHandler)

	// Serve static assets (CSS, images)
	http.Handle("/frontend/", http.StripPrefix("/frontend/", http.FileServer(http.Dir("frontend"))))

	// Start server
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		logHistory(fmt.Sprintf("Server stopped due to error: %s", err))
		fmt.Fprintln(os.Stderr, "Server error:", err)
	}
}
