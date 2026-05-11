package routes

import (
	"encoding/json"
	"log"
	"net/http"

	"sgougoupractice/fetch"
	"sgougoupractice/handlers"
	"sgougoupractice/helpers"
)

// Serve the concert dates page
func serveDatesPage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "assets/static/dates.html")
}

// Fetch and serve concert dates for a specific artist based on their ID
func serveDatesData(w http.ResponseWriter, r *http.Request) error {

	ok, id := helpers.CheckUrl(r.URL.Path)
	if !ok {

		return &handlers.HTTPError{
			Status:  http.StatusBadRequest,
			Message: "bad-format or invalid ",
		}
	}

	// Fetch the list of artists
	artists, err := fetch.FetchArtists()
	if err != nil {
		log.Println("Error fetching artists:", err)
		//handlers.ErrorHandler(w, r, http.StatusInternalServerError, "Failed to fetch artists.")
		return &handlers.HTTPError{
			Status:  http.StatusInternalServerError,
			Message: "unable to get artists",
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
			Message: "Invalid artist id",
		}
	}

	// Fetch concert dates associated with the artist
	dates, err := fetch.FetchDates(id)
	if err != nil {
		log.Println("Error fetching dates data:", err)

		return &handlers.HTTPError{
			Status:  http.StatusInternalServerError,
			Message: "error feching dates",
		}
	}

	// Structure the response with artist name and concert dates
	response := struct {
		ArtistName string   `json:"artistName"`
		Dates      []string `json:"dates"`
	}{
		ArtistName: artistName,
		Dates:      dates.Dates,
	}

	// Send the JSON response back to the client
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Println("Error encoding response:", err)

		return &handlers.HTTPError{
			Status:  http.StatusInternalServerError,
			Message: "encoding-error",
		}
	}
	return nil
}
