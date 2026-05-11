package main

import (
	"crypto/tls"
	"forum-authentication/internal/frontend/app"
	"log"
	"net/http"
	"time"
)

func main() {
	// Frontend app
	handler, CONFIG := app.New()

	// TLS
	cert, err := tls.LoadX509KeyPair(CONFIG.Tls.Certification, CONFIG.Tls.Key)
	if err != nil {
		log.Fatalf("failed to load cert/key: %v", err)
	}
	tlsCfg := &tls.Config{Certificates: []tls.Certificate{cert}}

	// Routes

	srv := &http.Server{
		Addr:         ":3000",
		Handler:      handler,
		TLSConfig:    tlsCfg,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Println("client listening on https://localhost:3000")
	log.Fatal(srv.ListenAndServeTLS("", ""))
}
