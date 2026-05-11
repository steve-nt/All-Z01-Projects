package routes

import (
	"net/http"

	"sgougoupractice/handlers"
)

// function for a custom 404 http error page
func NotFoundWrapper(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler, pattern := h.(*http.ServeMux).Handler(r)
		if pattern == "" {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("<h1>404 - Page Not Found</h1>"))
			return
		}
		handler.ServeHTTP(w, r)
	})
}

func RouteHandler(h *Handler) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			handlers.CustomNotFound(w, r)
			return
		}
		serveHome(w, r)
	})

	mux.Handle("/artists", handlers.Apphandler(serveArtists))
	mux.Handle("/static/", http.StripPrefix("/static/", handlers.StaticHandler()))
	//mux.Handle("/assets/images/", http.StripPrefix("/assets/images/", http.FileServer(http.Dir("assets/images"))))

	mux.Handle("/api/locations/", handlers.Apphandler(h.serveLocations))
	mux.Handle("/api/coords", handlers.Apphandler(h.ServeCoords))
	mux.Handle("/api/locations-list", handlers.Apphandler(serveAllLocations))
	mux.Handle("/api/dates/", handlers.Apphandler(serveDatesData))
	mux.Handle("/api/relations/", handlers.Apphandler(serveRelationsData))

	mux.HandleFunc("/dates.html", serveDatesPage)
	mux.HandleFunc("/relations.html", serveRelationsPage)
	mux.HandleFunc("/locations.html", serveLocationsPage)

	mux.Handle("/suggestions", handlers.Apphandler(handlers.SuggestionsHandler))

	// Ensure this route matches the fetch in your JS
	mux.Handle("/filterArtists", handlers.Apphandler(serveFilterArtists))

	mux.HandleFunc("/map.html", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "assets/static/map.html")
	})

	return mux
}
