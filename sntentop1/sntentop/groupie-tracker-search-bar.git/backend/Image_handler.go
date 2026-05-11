package backend

import (
	"net/http"
	"os"
)

// ImageHandler is a function that handles HTTP requests related to images.
func ImageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		files, err := os.ReadDir("frontend/images")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			http.ServeFile(w, r, "templates/500.html")
			return
		}
		// Iterate over the list of files in the directory
		for _, file := range files {
			if r.URL.Path == "/frontend/images/"+file.Name() {
				fs := http.Dir("frontend/images")
				http.StripPrefix("/frontend/images", http.FileServer(fs)).ServeHTTP(w, r)
				return
			}
		}
		// If no matching file is found, set the HTTP status code to 404 (Not Found)
		w.WriteHeader(http.StatusNotFound)
		http.Redirect(w, r, "/404", http.StatusFound)
		return
	} else {
		// If the HTTP method is not GET, set the HTTP status code to 405 (Method Not Allowed)
		w.WriteHeader(http.StatusMethodNotAllowed)
		http.ServeFile(w, r, "templates/405.html")
		return
	}
}
