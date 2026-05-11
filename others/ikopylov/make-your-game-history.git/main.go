package main

import (
	"log"
	"net/http"

	"web/internal/scoreboard"
	"web/internal/server"
)

func main() {
	store, err := scoreboard.NewStore("scores.json")
	if err != nil {
		log.Fatalf("failed to load score store: %v", err)
	}
	svc := scoreboard.NewService(store)
	handler := server.NewHandler(svc)

	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)
	mux.Handle("/", http.FileServer(http.Dir("./")))

	if err := server.Run(server.Config{Addr: ":8080"}, mux); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
