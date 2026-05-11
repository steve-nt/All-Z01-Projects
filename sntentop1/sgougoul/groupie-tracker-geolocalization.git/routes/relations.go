package routes

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"sgougoupractice/fetch"
	"sgougoupractice/handlers"
	"sgougoupractice/helpers"
)

// Serve the relations page
func serveRelationsPage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "assets/static/relations.html")
}

// Fetch and serve relations data for an artist
func serveRelationsData(w http.ResponseWriter, r *http.Request) error {
	ok, id := helpers.CheckUrl(r.URL.Path)
	if !ok {
		return &handlers.HTTPError{
			Status:  http.StatusBadRequest,
			Message: "bad-format or invalid ",
		}
	}

	// Fetch the list of artists from the external API
	artists, err := fetch.FetchArtists()
	if err != nil {
		log.Println("Error fetching artists:", err)
		return &handlers.HTTPError{
			Status:  http.StatusInternalServerError,
			Message: "error feching artists",
		}
	}

	// Find the artist's name based on the provided ID
	var artistName string
	for _, artist := range artists {
		if artist.ID == id {
			artistName = artist.Name
			break
		}
	}
	// Return an error if the artist is not found
	if artistName == "" {
		return &handlers.HTTPError{
			Status:  http.StatusBadRequest,
			Message: "Invalid asrtist",
		}
	}

	// Fetch the artist's relations data from the external API
	relation, err := fetch.FetchRelations(id)
	if err != nil {
		log.Println("Error fetching relations data:", err)
		return &handlers.HTTPError{
			Status:  http.StatusInternalServerError,
			Message: "error fetcing relations",
		}
	}

	// Convert the relations data (DatesLocations map) into a readable list
	var relationDetails []string
	for key, related := range relation.DatesLocations {
		relationDetails = append(relationDetails, fmt.Sprintf("%s: %v", key, related))
	}

	// Construct the JSON response with artist name and relations
	response := struct {
		ArtistName      string   `json:"artistName"`
		RelationDetails []string `json:"relationDetails"`
	}{
		ArtistName:      artistName,
		RelationDetails: relationDetails, // List of relations from the map
	}

	// Send the JSON response back to the client
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Println("Error encoding response:", err)
		return &handlers.HTTPError{
			Status:  http.StatusInternalServerError,
			Message: "encoding error",
		}
	}
	return nil
}
