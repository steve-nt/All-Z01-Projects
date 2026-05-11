package handlers

import "net/http"

func RegisterRoutes(mux *http.ServeMux) {

	// Static files
	fs := http.FileServer(http.Dir("web"))
	mux.Handle("/web/", http.StripPrefix("/web/", fs))

	// API routes
	mux.HandleFunc("/api/maxmembers", HandleMaxMembers)
	mux.HandleFunc("/api/minmax", HandleMinMax)
	mux.HandleFunc("/geolocations", GeocodeAddressHandler)
	mux.HandleFunc("/all_locations", HandleAllLocations)

	// Page routes
	mux.HandleFunc("/", HandleHome)
	mux.HandleFunc("/about", HandleAbout)
	mux.HandleFunc("/index", HandleIndex)
	mux.HandleFunc("/Artist/", HandlePage)
	mux.HandleFunc("/404", ErrorHandler)

	// Functional routes
	mux.HandleFunc("/search", HandleSearch)
	mux.HandleFunc("/filters", HandleFilters)

}
