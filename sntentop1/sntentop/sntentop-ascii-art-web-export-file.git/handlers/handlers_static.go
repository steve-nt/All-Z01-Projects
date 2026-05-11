package handlers

import (
	"log"      // Provides simple logging functions like Println, Fatalf, etc.
	"net/http" // Implements HTTP client and server functionalities
	"os"       // Provides functions to interact with the operating system (e.g., file operations)
)

// ServeStaticFiles is a function that configures and registers a handler to serve static files.
func ServeStaticFiles() {
	// Define the directory where static files are stored.
	staticDir := "static"
	// Check if the static directory exists using os.Stat, which returns file info or an error.
	// os.IsNotExist checks if the error returned by os.Stat is due to the file/directory not existing.
	if _, err := os.Stat(staticDir); os.IsNotExist(err) {
		// If the static directory does not exist, log an error and exit the program using log.Fatalf.
		log.Fatalf("Static directory does not exist: %s", staticDir)
	}
	// Log a success message if the directory exists.
	log.Println("Static directory cofiguration succeded")
	// Log a message indicating that a static file handler is being registered.
	log.Println("Registering static file handler")
	// Register a handler for HTTP requests with paths starting with "/static/".
	http.Handle("/static/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Log the path of the requested static file.
		log.Println("Static file requested:", r.URL.Path)
		// Serve the static files.
		// http.StripPrefix removes the "/static/" prefix from the request path before serving the file.
		// http.FileServer serves files from the given directory (in this case, the "static" directory).
		http.StripPrefix("/static/", http.FileServer(http.Dir("static"))).ServeHTTP(w, r)
	}))
}
