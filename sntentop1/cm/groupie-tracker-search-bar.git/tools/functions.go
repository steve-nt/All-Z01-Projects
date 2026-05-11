package tools

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

var artists []Artist

// Fetch data from the API and parse it
func fetchData(url string, target interface{}) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to fetch data: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %v", err)
	}

	err = json.Unmarshal(body, target)
	if err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %v", err)
	}
	return nil
}

// Load and link data
func loadData() {
	// Fetch all ARTISTS
	err := fetchData("https://groupietrackers.herokuapp.com/api/artists", &artists)
	if err != nil {
		log.Printf("Error fetching artists: %v\n", err)
		artists = []Artist{}
	}

	// Iterate over the artists and fetch the RELATIONS
	for i := range artists {
		artist := &artists[i]

		var relationData RelationData
		err := fetchData(fmt.Sprintf("https://groupietrackers.herokuapp.com/api/relation/%d", artist.ID), &relationData)
		if err != nil {
			log.Printf("Error fetching relations for artist %d: %v\n", artist.ID, err)
			continue
		}

		// Extract locations from the datesLocations field
		for location := range relationData.DatesLocations {
			artist.Locations = append(artist.Locations, location)
		}

		// Now you have the artist with the extracted locations
		artist.Relations = relationData.DatesLocations
	}
}

func Error404Page(w http.ResponseWriter, r *http.Request) {
	// Render the "error404.html" template
	err := tmpl.ExecuteTemplate(w, "error404.html", nil)
	if err != nil {
		http.Error(w, "Error-Internal 500", http.StatusInternalServerError)
	}
}

func ListenAndServe() {
	loadData()

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/artist", artistDetailsHandler)
	http.HandleFunc("/search", searchHandler)
	http.HandleFunc("/searchSuggestions", searchSuggestionsHandler)
	http.HandleFunc("/artist/locations", artistLocationsHandler)

	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css"))))

	fmt.Println("Server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
