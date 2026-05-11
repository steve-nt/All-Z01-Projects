package handlers

import (
	"encoding/json"
	"gtracker/api"
	"gtracker/internal/services"
	"log"
	"net/http"
	"strings"
)

// GeolocationHandler handles the request to fetch geolocation data
func GeolocationHandler(w http.ResponseWriter, r *http.Request) {
	currArtist := r.URL.Query().Get("concert")
	var concerts, locations, dates, latitudes, longitudes []string
	allArtists, err := services.GetArtists()
	if err != nil {
		http.Error(w, "Failed to fetch artists", http.StatusInternalServerError)
		log.Printf("Error fetching artists: %v", err)
		return
	}
	for _, artist := range allArtists {
		if currArtist == artist.Name {
			concerts = artist.ConcertsFormatted
		}
	}

	if concerts == nil {
		http.Error(w, "Concert information is required", http.StatusBadRequest)
		return
	}

	// Collect locations from the concerts
	for _, concert := range concerts {
		tempConcert := strings.Split(concert, ":")
		if len(tempConcert) > 0 {
			locations = append(locations, tempConcert[0]) // Get the location part before ":"
			dates = append(dates, tempConcert[1])
		}
	}

	for _, location := range locations {

		lat, lon, err := api.GeocodeAddress(location)
		if err != nil {
			return
		}
		latitudes = append(latitudes, lat)
		longitudes = append(longitudes, lon)

	}

	// Return the results (location, latitude, longitude) as JSON
	response := map[string]any{
		"artist":     currArtist,
		"dates":      dates,
		"locations":  locations,
		"latitudes":  latitudes,
		"longitudes": longitudes,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
