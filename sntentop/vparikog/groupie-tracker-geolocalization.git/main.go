package main

import (
	groupie_tracker_search "groupie_tracker_search/funcs"

	"log"
	"net/http"
)

func main() {

	mux := http.NewServeMux()

	staticDir := http.Dir("css")
	mux.Handle("/css/", http.StripPrefix("/css/", http.FileServer(staticDir)))

	staticJSDir := http.Dir("JavaScript")
	mux.Handle("/JavaScript/", http.StripPrefix("/JavaScript/", http.FileServer(staticJSDir)))

	mux.HandleFunc("/api/locations", groupie_tracker_search.HandleLocations)

	mux.HandleFunc("/", groupie_tracker_search.MainPage)
	// Artist details route
	mux.HandleFunc("/artist", groupie_tracker_search.ArtistPage)

	mux.HandleFunc("/search", groupie_tracker_search.SearchHandler)

	mux.HandleFunc("/coordinates", groupie_tracker_search.CoordinatesHandler)

	log.Println("Server is running at http://localhost:2000")
	log.Fatal(http.ListenAndServe(":2000", mux))

}
