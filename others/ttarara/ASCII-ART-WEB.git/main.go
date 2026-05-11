package main

import (
	"log"
	"net/http"
	"stylizeWeb/handlers"
)

func main() {
	handlers.ServeStaticFiles()

	mux := http.NewServeMux()

	mux.HandleFunc("/home", handlers.PgHome)
	mux.HandleFunc("/project", handlers.PgProject)
	mux.HandleFunc("/converter", handlers.PgConverter)
	mux.HandleFunc("/team", handlers.PgTeam)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			http.Redirect(w, r, "/home", http.StatusFound) // Redirect to the home page
			return
		}
		if _, pattern := mux.Handler(r); pattern == "" { // Check if the request matches a registered route
			log.Println("Error 404: Page not found")
			handlers.NotFoundHandler(w, r) // Use custom 404 handler for unmatched routes
			return
		}
		mux.ServeHTTP(w, r) // Pass the request to the matched handler
	})

	log.Println("Starting server on http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
