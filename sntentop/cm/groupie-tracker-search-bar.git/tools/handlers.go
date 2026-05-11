package tools

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

// we read and parse all .html files in the templates directory
// and we combine them into a single *template.Template object stored in tmpl variable
var tmpl = template.Must(template.ParseGlob("templates/*.html"))

// Handlers for the web server
func homeHandler(w http.ResponseWriter, r *http.Request) {
	// handling Error 404:
	if r.URL.Path != "/" {
		Error404Page(w, r)
		return
	}

	// Render the filtered artists and load the HTML template
	tmplData := filters(r)
	err := tmpl.ExecuteTemplate(w, "index.html", tmplData)
	if err != nil {
		http.Error(w, "HTTP status 500 - Could not render template", http.StatusInternalServerError)
		return
	}
}
func searchHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query")
	if query == "" {
		// Redirect to the home page ("/")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	// Filter artists based on the search query
	var results []Artist
	for _, artist := range artists {
		// Check if the query matches the name or any location
		if containsIgnoreCase(artist.Name, query) || containsInSliceIgnoreCase(artist.Locations, query) || containsInSliceIgnoreCase(artist.Members, query) || containsIgnoreCase(artist.FirstAlbum, query) || containsInt(artist.CreationDate, query) {
			results = append(results, artist)
		}
	}
	// Serve results
	tmpl, err := template.ParseFiles("templates/search_results.html")
	if err != nil {
		http.Error(w, "HTTP status 500 - Could not load template", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, results)
}
func searchSuggestionsHandler(w http.ResponseWriter, r *http.Request) {
	query := strings.TrimSpace(r.URL.Query().Get("query"))
	queryLower := strings.ToLower(query)

	var suggestions []map[string]string
	locationMap := make(map[string]struct{}) // Track unique locations

	for _, artist := range artists {
		creationDateStr := strconv.Itoa(artist.CreationDate)

		if strings.HasPrefix(strings.ToLower(artist.Name), queryLower) {
			suggestions = append(suggestions, map[string]string{
				"text": artist.Name,
				"type": "- Artist/Band",
			})
		}

		if strings.HasPrefix(strings.ToLower(artist.FirstAlbum), queryLower) {
			suggestions = append(suggestions, map[string]string{
				"text": artist.FirstAlbum,
				"type": "- First Album of " + artist.Name,
			})
		}

		for _, member := range artist.Members {
			if strings.HasPrefix(strings.ToLower(member), queryLower) {
				suggestions = append(suggestions, map[string]string{
					"text": member,
					"type": "- Member of " + artist.Name,
				})
			}
		}

		for location := range artist.Relations { // Assuming artist.Relations is a map
			if containsIgnoreCase(location, queryLower) {
				if _, exists := locationMap[location]; !exists { // Check if already added
					suggestions = append(suggestions, map[string]string{
						"text": location,
						"type": "- Location",
					})
					locationMap[location] = struct{}{} // Mark as added
				}
			}
		}

		if strings.Contains(creationDateStr, queryLower) {
			suggestions = append(suggestions, map[string]string{
				"text": creationDateStr,
				"type": "- Creation Date of" + artist.Name,
			})
		}
	}
	//IDEA: sort the map before passing it to js!!!!
	sort.Slice(suggestions, func(i, j int) bool {
		aStarts := strings.HasPrefix(strings.ToLower(suggestions[i]["text"]), queryLower)
		bStarts := strings.HasPrefix(strings.ToLower(suggestions[j]["text"]), queryLower)

		if aStarts && !bStarts {
			return true // `a` comes first
		}
		if !aStarts && bStarts {
			return false // `b` comes first
		}
		return strings.ToLower(suggestions[i]["text"]) < strings.ToLower(suggestions[j]["text"]) // Alphabetical as tiebreaker
	})

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(suggestions); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}

func containsInSliceIgnoreCase(slice []string, query string) bool {
	query = strings.ToLower(query)
	for _, item := range slice {
		if strings.Contains(strings.ToLower(item), query) {
			return true
		}
	}
	return false
}
func containsIgnoreCase(str, substr string) bool {
	return strings.Contains(strings.ToLower(str), strings.ToLower(substr))
}
func containsInt(value int, substr string) bool {
	str := strconv.Itoa(value)
	return strings.HasPrefix(strings.ToLower(str), strings.ToLower(substr))
}
func artistLocationsHandler(w http.ResponseWriter, r *http.Request) {
	artistID := r.URL.Query().Get("id")

	var locationsMap = make(map[string]struct{}) // Use map to remove duplicates
	var found bool

	for _, artist := range artists {
		if strconv.Itoa(artist.ID) == artistID {
			for _, location := range artist.Locations {
				locationsMap[location] = struct{}{} // Store unique locations
			}
			found = true
			break
		}
	}

	if !found {
		http.Error(w, "Artist not found", http.StatusNotFound)
		return
	}

	// Convert map keys back to slice
	locations := make([]string, 0, len(locationsMap))
	for location := range locationsMap {
		locations = append(locations, location)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(locations)
}

func artistDetailsHandler(w http.ResponseWriter, r *http.Request) {
	// Extract artist ID from query
	id := r.URL.Query().Get("id")
	for _, artist := range artists {
		if fmt.Sprintf("%d", artist.ID) == id {
			// Render the "artist_details.html" template
			err := tmpl.ExecuteTemplate(w, "artist_details.html", artist)
			if err != nil {
				http.Error(w, "HTTP status 500 - Could not render template", http.StatusInternalServerError)
				return
			}
			return
		}
	}

	// Handling Error 404:
	Error404Page(w, r)
}
