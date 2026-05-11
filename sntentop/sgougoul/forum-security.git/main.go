package main

import (
	"log"

	"forum/db"
	"forum/handlers"
	"forum/server"
)

// main bootstraps the application.
// AUDIT: keeps initialization separate from server logic.
func main() {

	db.Init()

	if err := handlers.InitTemplates(); err != nil {
		log.Fatal("Error parsing templates:", err)
	}

	srv := server.NewServer()

	log.Fatal(srv.Run())
}