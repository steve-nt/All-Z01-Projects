package main

import (
	"log"
	"net/http"

	"forum-authentication/internal/backend/app"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	handler, db := app.New()
	defer db.Close()

	server := &http.Server{
		Addr:    ":8080",
		Handler: handler,
	}

	log.Println("Listening on :8080")
	log.Fatal(server.ListenAndServe())
}
