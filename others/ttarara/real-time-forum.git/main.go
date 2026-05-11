// main.go
// Responsibilities: bootstraps the application by initializing the database, wiring routes, and starting the HTTP server listener.
package main

import (
	"fmt"
	"log"
	"net/http"

	"realtimeforum/internals/database"
	"realtimeforum/internals/handlers"
)

func main() {
	// Initialize database connections and migrations before serving requests.
	database.InitializeDatabase()

	fmt.Println("Server running on http://localhost:8081")

	// Setup routes (uses handlers package to register HTTP endpoints).
	handlers.SetupRoutes()

	// Start server
	if err := http.ListenAndServe(":8081", nil); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
