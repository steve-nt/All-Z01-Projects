package router

import (
	"gtracker/internal/handlers"
	"net/http"
)

func InitRoutes() {

	fs := http.FileServer(http.Dir("web/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", handlers.ArtistsHandler)
	http.HandleFunc("/search", handlers.SearchHandler)
	http.HandleFunc("/suggestions", handlers.SuggestionsHandler)
	http.HandleFunc("/geolocations", handlers.GeolocationHandler)
	http.HandleFunc("/filter", handlers.FilterArtistsHandler)
	http.HandleFunc("/locations", handlers.FetchAvailableConcertLocations)

}
