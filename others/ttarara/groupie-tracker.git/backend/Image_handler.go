package backend
import (
	"net/http"
	"os"
)

func ImageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		files, err := os.ReadDir("frontend/images")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			http.ServeFile(w, r, "templates/500.html")
			return
		}
		for _, file := range files {
			if r.URL.Path == "/frontend/images/"+file.Name() {
				fs := http.Dir("frontend/images")
				http.StripPrefix("/frontend/images", http.FileServer(fs)).ServeHTTP(w, r)
				return
			}
		}
		w.WriteHeader(http.StatusNotFound)
		http.Redirect(w, r, "/404", 302)
		return
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		http.ServeFile(w, r, "templates/405.html")
		return
	}
}
