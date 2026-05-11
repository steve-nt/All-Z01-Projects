package main

import (
	"log"                 // Standard library for logging messages to the console.
	"net/http"            // Standard library for HTTP client and server implementations.
	"stylizeWeb/handlers" // Custom package that contains functions for handling HTTP requests.
)

func main() {
	// Call a function from the "handlers" package to serve static files like CSS, JavaScript, or images
	handlers.ServeStaticFiles()

	// Create a new ServeMux (multiplexer) to manage URL routes and their associated handlers.
	mux := http.NewServeMux()

	// Register a handler for the "/home" route using a function from the "handlers" package.
	mux.HandleFunc("/home", handlers.PgHome)
	// Register a handler for the "/project" route using a function from the "handlers" package.
	mux.HandleFunc("/project", handlers.PgProject)
	// Register a handler for the "/converter" route using a function from the "handlers" package.
	mux.HandleFunc("/converter", handlers.PgConverter)
	// Register a handler for the "/team" route using a function from the "handlers" package.
	mux.HandleFunc("/team", handlers.PgTeam)
	// Register a handler for the "/download" route using a function from the "handlers" package
	mux.HandleFunc("/download", handlers.HandleDownload)
	// Register a handler for the "/check-file" route using a function from the "handlers" package.
	mux.HandleFunc("/check-file", handlers.HandleCheckFile)
	// Register a handler for the "/export-zip" route using a function from the "handlers" package
	mux.HandleFunc("/export-zip", handlers.HandleExportZip)

	// Set up a default handler for all incoming requests.
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Redirect the root path ("/") to "/home".
		if r.URL.Path == "/" {
			http.Redirect(w, r, "/home", http.StatusFound) // Redirect to the home page
			return
		}
		// Check if the request matches any of the registered routes in the mux.
		if _, pattern := mux.Handler(r); pattern == "" {
			// Log an error message if the page is not found.
			log.Println("Error 404: Page not found")
			// Call a custom 404 handler from the "handlers" package.
			handlers.NotFoundHandler(w, r)
			return
		}

		// If a route is matched, serve the request using the corresponding handler.
		mux.ServeHTTP(w, r)
	})
	// Log a message to indicate the server is starting on port 8080.
	log.Println("Starting server on http://localhost:8080")
	// Start the HTTP server on port 8080. If an error occurs, log it and exit.
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
