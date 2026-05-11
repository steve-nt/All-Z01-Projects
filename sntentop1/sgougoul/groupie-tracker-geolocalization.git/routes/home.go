package routes

import (
	"net/http"
)

// Serve the home page
func serveHome(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "assets/static/home.html")
}

// Function to deliberately trigger a panic (for debugging)
func servePanic(w http.ResponseWriter, r *http.Request) {
	panic("This is a test panic!")
}
