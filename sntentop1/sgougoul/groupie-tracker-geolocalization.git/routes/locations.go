package routes

import (
	"encoding/json"
	"log"
	"net/http"
	"regexp"
	"sort"
	"strconv"
	"sync"

	"sgougoupractice/fetch"
	"sgougoupractice/handlers"

	fetcher "sgougoupractice/geocoding/fetch"
	geocode "sgougoupractice/geocoding/service"
)

// Handler struct contains dependencies for fetching artist data, geocoding, and caching location data.
type Handler struct {
	Fetcher     fetcher.Fetcher
	GeoSvc      *geocode.GeocodingService
	RawLocCache map[int][]string
	rawLocMux   sync.RWMutex
}

// NewHandler is a constructor function for creating a new Handler instance.
func NewHandler(f fetcher.Fetcher, g *geocode.GeocodingService) *Handler {
	return &Handler{
		Fetcher:     f,
		GeoSvc:      g,
		RawLocCache: make(map[int][]string),
	}
}

var reg = regexp.MustCompile(`^\d+$`)

// serveLocations handles requests to fetch the locations of an artist based on their ID from the query parameter.
// It validates the artist ID, fetches artist data, fetches location data, stores it in cache, and sends the response back.
func (h *Handler) serveLocations(w http.ResponseWriter, r *http.Request) error {
	// Retrieve artist ID from query parameters and validate
	idq := r.URL.Query().Get("id")

	if reg.FindString(idq) == "" {
		return &handlers.HTTPError{
			Status:  http.StatusBadRequest,
			Message: "Bad Request",
		}
	}
	id, err := strconv.Atoi(idq)
	if err != nil {
		return &handlers.HTTPError{
			Status:  http.StatusBadRequest,
			Message: "Bad Request",
		}
	}
	// Fetch artists data and find the artist by ID
	artists, err := fetch.FetchArtists()
	if err != nil {
		log.Println("Error fetching artists:", err)
		return &handlers.HTTPError{
			Status:  http.StatusInternalServerError,
			Message: "error getting artists ",
		}
	}
	// Find artist by ID
	var artistName string
	for _, artist := range artists {
		if artist.ID == id {
			artistName = artist.Name
			break
		}
	}
	if artistName == "" && id > 0 {
		return &handlers.HTTPError{
			Status:  http.StatusBadRequest,
			Message: "not such artists ",
		}
	}
	// Fetch location data for the artist
	location, err := h.Fetcher.Fetch(r.Context(), id)
	if err != nil {
		log.Println("Error fetching location data:", err)
		handlers.ErrorHandler(w, r, http.StatusInternalServerError, "Failed to fetch location data.")
		return &handlers.HTTPError{
			Status:  http.StatusInternalServerError,
			Message: "error getting locations ",
		}
	}
	// Store the fetched locations in cache with thread safety
	h.rawLocMux.Lock()
	h.RawLocCache[id] = location.Locations
	h.rawLocMux.Unlock()
	log.Printf("[serveLocations] storing %d locations for artist %d", len(location.Locations), id)
	response := struct {
		ArtistName string   `json:"artistName"`
		Locations  []string `json:"locations"`
	}{
		ArtistName: artistName,
		Locations:  location.Locations,
	}

	// Send the response back to the client
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Println("Error encoding response:", err)
		return &handlers.HTTPError{
			Status:  http.StatusInternalServerError,
			Message: "error encoding ",
		}
	}
	return nil
}

// serveLocationsPage serves the static HTML page for the location interface.
func serveLocationsPage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "assets/static/locations.html")
}

// serveAllLocations fetches all the unique locations for all artists and returns them as a sorted list.
func serveAllLocations(w http.ResponseWriter, r *http.Request) error {
	// Fetch all artists
	artists, err := fetch.FetchArtists()
	if err != nil {
		log.Println("Error fetching artists:", err)
		return &handlers.HTTPError{
			Status:  http.StatusInternalServerError,
			Message: "Failed to fetch artists",
		}
	}
	// Set to store unique locations
	locationSet := make(map[string]struct{})
	// Fetch and accumulate locations for each artist
	for _, artist := range artists {
		locData, err := fetch.FetchLocations(artist.ID)
		if err != nil {
			log.Printf("Skipping locations for artist ID %d due to error: %v", artist.ID, err)
			continue
		}
		// Add each location to the set (duplicates will be ignored)
		for _, loc := range locData.Locations {
			locationSet[loc] = struct{}{}
		}
	}
	// Convert the set of locations into a sorted slice
	uniqueLocations := make([]string, 0, len(locationSet))
	for loc := range locationSet {
		uniqueLocations = append(uniqueLocations, loc)
	}
	sort.Strings(uniqueLocations)
	// Prepare the response and encode it as JSON
	resp := struct {
		Locations []string `json:"locations"`
	}{
		Locations: uniqueLocations,
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		log.Println("Error encoding locations response:", err)
		return &handlers.HTTPError{
			Status:  http.StatusInternalServerError,
			Message: "Failed to encode locations",
		}
	}
	return nil
}

// ServeCoords handles requests to fetch geographical coordinates for locations of an artist.
// It uses the cache for locations and the geocoding service to fetch the coordinates.
func (h *Handler) ServeCoords(w http.ResponseWriter, r *http.Request) error {
	// Retrieve artist ID from query parameters and validate
	idq := r.URL.Query().Get("id")
	if reg.FindString(idq) == "" {
		return &handlers.HTTPError{Status: http.StatusBadRequest, Message: "Bad Request"}
	}
	id, _ := strconv.Atoi(idq)
	// Check if the locations for the artist are already cached
	h.rawLocMux.RLock()
	slugs, ok := h.RawLocCache[id]
	h.rawLocMux.RUnlock()
	if !ok {

		return &handlers.HTTPError{Status: http.StatusBadRequest, Message: "Must fetch locations first"}
	}
	log.Printf("[ServeCoords] looking up coords for artist %d", id)
	// Use the geocoding service to fetch coordinates for the cached locations
	coords, err := h.GeoSvc.BatchGeocode(r.Context(), slugs)
	if err != nil {
		log.Println("Error geocoding locations:", err)
		return &handlers.HTTPError{Status: http.StatusInternalServerError, Message: "error getting coords"}
	}
	// Prepare a map of locations to coordinates
	coordsMap := make(map[string]geocode.LocationCoord, len(coords))
	for _, c := range coords {
		coordsMap[c.Name] = c
	}
	// Send the geocoded coordinates as a JSON response
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(coordsMap)
}
