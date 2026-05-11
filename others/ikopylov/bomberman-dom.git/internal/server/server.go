package server

import (
	"log"
	"net/http"
)

type Config struct {
	Addr string
}

func Run(cfg Config, mux *http.ServeMux) error {
	log.Printf("Server running on http://localhost%s\n", cfg.Addr)
	return http.ListenAndServe(cfg.Addr, mux)
}
