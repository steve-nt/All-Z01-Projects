package groupie_tracker_search

import (
	"encoding/json"
	"net/http"
	"strconv"
)

func CoordinatesHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid artist ID", http.StatusBadRequest)
		return
	}

	artists, err := getCachedArtists()
	if err != nil {
		http.Error(w, "Failed to fetch artists", http.StatusInternalServerError)
		return
	}

	var selectedArtist *Artist
	for i := range artists {
		if artists[i].ID == id {
			selectedArtist = &artists[i]
			// Geolocate here just like in ArtistPage
			if len(selectedArtist.GeoCoordinates) == 0 {
				geolocated := GeolocateLocations([]Artist{*selectedArtist})
				*selectedArtist = geolocated[0]
			}
			break
		}
	}

	if selectedArtist == nil {
		http.Error(w, "Artist not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(selectedArtist.GeoCoordinates)
	if err != nil {
		http.Error(w, "Failed to encode coordinates", http.StatusInternalServerError)
	}
}
