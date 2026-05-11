package entry

import (
	"fmt"
	"groupie/geocoding"
	api "groupie/handlers"
	"log"
	"net/http"
)

func Start() {

	var err error
	api.Artists, api.RelationData, err = api.LoadArtistData()
	if err != nil {
		log.Fatal("Critical error on init: ", err.Error())
	}

	http.Handle("/templates/", http.StripPrefix("/templates/", http.FileServer(http.Dir("../templates"))))
	// http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("../images"))))

	http.HandleFunc("/", api.MainHandler)
	http.HandleFunc("/search", api.SuggestionsHandler)
	http.HandleFunc("/map", api.MapHandler)
	http.HandleFunc("/markerHandler", api.MarkerHandler)

	err = geocoding.LoadGeocodeData()
	if err != nil {
		fmt.Println("ERROR: failed to load geocoding data:", err)
	}
	go geocoding.GeocodeLogger()
	go geocoding.GeocodeDownloader()

	log.Println("Server running on :8080")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("ERROR: ", err)
	}
}
