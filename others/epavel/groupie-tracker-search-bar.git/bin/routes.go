package bin

import (
	"net/http"
)

// RegisterRoutes registers the routes for the application
func RegisterRoutes() *http.ServeMux {
	router := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("public"))
	router.Handle("/public/", http.StripPrefix("/public", noCacheMiddleware(fileServer)))

	router.Handle("/filter/data", http.HandlerFunc(HandleFilterData))
	router.Handle("/artist/", http.HandlerFunc(handleIndividualArtistRequest))
	router.Handle("/", Middleware(http.HandlerFunc(handleHome)))

	return router
}
