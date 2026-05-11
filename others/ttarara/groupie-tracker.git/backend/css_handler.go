package backend

import (
	"net/http"
	"os"
)

func CssHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		files, err := os.ReadDir("frontend/css")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			http.ServeFile(w, r, "templates/500.html")
			return
		}
		for _, file := range files {
			if r.URL.Path == "/frontend/css/"+file.Name() {
				fs := http.Dir("frontend/css")
				http.StripPrefix("/frontend/css", http.FileServer(fs)).ServeHTTP(w, r)
				return
			}
		}
		http.ServeFile(w, r, "templates/404.html")
		return
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		http.ServeFile(w, r, "templates/405.html")
		return
	}
}
