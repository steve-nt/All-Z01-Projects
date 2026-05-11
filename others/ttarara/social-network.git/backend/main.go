package main

import (
	"log"
	"net/http"
	"social-network/backend/pkg/db/sqlite"
	"social-network/backend/utils"

	// Import handlers package (authentication folder uses package name "handlers")
	handlers "social-network/backend/pkg/authentication"
	"social-network/backend/pkg/groups"
	"social-network/backend/pkg/messages"
	"social-network/backend/pkg/notifications"
	"social-network/backend/pkg/posts"
	"social-network/backend/pkg/websocket"
)

func main() {
	// Step 1: Initialize the database
	// This will create the database file, open a connection, and apply all migrations
	// Database location: data/social_network.db (relative to where server runs)
	//
	// IMPORTANT: Run the server from the backend/ directory:
	//   cd backend
	//   go run main.go
	//
	// The 'data/' directory will be created automatically if it doesn't exist
	dbPath := "data/social_network.db"

	if err := sqlite.InitDB(dbPath); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Ensure database connection is closed when server shuts down
	defer func() {
		if err := sqlite.CloseDB(); err != nil {
			log.Printf("Error closing database: %v", err)
		}
	}()

	// Step 2: Set up your HTTP routes

	// Home page
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		// Use the existing template service with auth context
		utils.FileServiceWithAuth("home.html", w, r, nil)
	})

	// Serve uploaded images and thumbnails
	http.Handle("/frontend/uploads/", http.StripPrefix("/frontend/uploads/", http.FileServer(http.Dir("frontend/uploads"))))

	// Serve WebSocket test page
	http.HandleFunc("/test_websocket.html", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "test_websocket.html")
	})

	// Health check endpoint (no middleware needed)
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		// Test database connection
		if err := sqlite.Ping(); err != nil {
			http.Error(w, "Database connection failed", http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK - Database connected"))
	})

	// Setup authentication routes (register, login, logout, etc.)
	// This includes middleware wrapping for logging and authentication
	// Note: SetupAuthRoutes is in the handlers package (authentication folder)
	handlers.SetupAuthRoutes()

	// Part 3: Posts & Groups
	posts.SetupPostRoutes()
	groups.SetupGroupRoutes()

	// Notifications API
	notifications.SetupNotificationRoutes()

	// Private Messaging
	messages.SetupMessageRoutes()

	// Part 4: WebSocket for real-time communication
	hub := websocket.NewHub()
	websocket.SetGlobalHub(hub) // Make hub accessible to other packages
	go hub.Run()                 // Start the hub in a goroutine
	websocket.SetupWebSocketRoutes(hub)

	// Step 3: Start the HTTP server
	port := ":8080"
	log.Printf("Server starting on port %s", port)
	log.Println("Database initialized and ready!")

	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
