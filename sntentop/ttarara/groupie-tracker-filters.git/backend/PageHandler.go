package backend

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"
)

// HandlePage serves the artist details page
func HandlePage(w http.ResponseWriter, r *http.Request) {
	// Ensure only GET requests are allowed; otherwise, return a 405 (Method Not Allowed) error
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		http.ServeFile(w, r, "templates/405.html")
		return
	}

	// Extract ID from URL path
	id := strings.TrimPrefix(r.URL.Path, "/Artist/")
	idd, err := strconv.Atoi(id)

	if err != nil || idd <= 0 || idd >= 53 {
		http.Redirect(w, r, "/404", http.StatusFound)
		return
	}

	// Construct the API URL to fetch artist data
	apiURL := "https://groupietrackers.herokuapp.com/api/artists/" + id

	// Fetch artist data using helper function
	artist, err := fetchData[Artist](apiURL)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	// Fetch extra details concurrently
	artist, err = fetchExtraDetails(artist)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	// Parse and execute the band details template
	tmpl, err := template.ParseFiles("templates/band.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		http.ServeFile(w, r, "templates/500.html")
		return
	}

	// Render the template with the artist data
	tmpl.Execute(w, artist)
}

// HandleIndex serves the index page with the artists
func HandleIndex(w http.ResponseWriter, r *http.Request) {
	LogHistory(fmt.Sprintf("Accessed Index Page - %s", r.RemoteAddr))

	// trigger error 500 by changing the link and commenting the func Init in main.go
	apiArtist := "https://groupietrackers.herokom/api/artists"


	// Fetch and store globally if not already populated
	if len(artists) == 0 {
		fetchedArtists, err := FetchArtists(apiArtist)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			http.ServeFile(w, r, "templates/500.html")
			return
		}
		artists = fetchedArtists // Assign to global variable
	}

	// Parse the index page template
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		http.ServeFile(w, r, "templates/500.html")
		return
	}

	// Render the index page with the artists
	tmpl.Execute(w, artists)
}

// serve the about page
func HandleAbout(w http.ResponseWriter, r *http.Request) {
	LogHistory(fmt.Sprintf("Accessed About Page - %s", r.RemoteAddr))
	http.ServeFile(w, r, "templates/about.html")
}
