package main

import (
	"ascii-art-web/handlers"
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	logger := log.New(os.Stdout, "APP: ", log.LstdFlags|log.Lshortfile)
	templates, err := handlers.InitializeTemplates("templates", logger)
	if err != nil {
		logger.Fatalf("Failed to initialize templates: %v", err)
	}
	handlers.SetTemplates(templates)
	logger.Println("Server started successfully")

	// Register routes
	http.HandleFunc("/", handlers.Index)

	http.HandleFunc("/ascii-art", handlers.Processor)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Start the server
	fmt.Println("HTTP SERVER RUNNING AT: http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Server failed to start: %v\n", err)
	}
}
