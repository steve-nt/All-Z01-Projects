package handlers

import (
	"log"
	"net/http"
	"os"
)

func ServeStaticFiles() {

	staticDir := "static"
	if _, err := os.Stat(staticDir); os.IsNotExist(err) {
		log.Fatalf("Static directory does not exist: %s", staticDir)
	}
	log.Println("Static directory cofiguration succeded")
	log.Println("Registering static file handler")
	http.Handle("/static/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Static file requested:", r.URL.Path)
		http.StripPrefix("/static/", http.FileServer(http.Dir("static"))).ServeHTTP(w, r)
	}))
}
