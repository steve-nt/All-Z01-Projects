package backend

import (
	"html/template"
	"net/http"
	"strconv"
	"strings"
)

// Handle Artist Page
func HandlePage(w http.ResponseWriter, r *http.Request) {
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

	// Fetch artist data concurrently
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

	// Parse and execute the template
	tmpl, err := template.ParseFiles("templates/band.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		http.ServeFile(w, r, "templates/500.html")
		return
	}
	tmpl.Execute(w, artist)
}
