package server

import (
	"log"
	"net/http"
	"os"
	"strings"

	"forum/handlers"
)

// Server represents the HTTP/HTTPS server instance.
// AUDIT: centralizes routing, middleware, and TLS configuration.
type Server struct {
	handler http.Handler
}

// NewServer builds the application server.
// AUDIT: separates server logic from main.go.
func NewServer() *Server {
	mux := http.NewServeMux()

	mux.Handle("/static/",
		http.StripPrefix("/static/",
			http.FileServer(http.Dir("static"))))

	registerRoutes(mux)

	// Middleware order:
	// 1. custom 404 handling
	// 2. security headers
	// 3. panic recovery
	//
	// AUDIT: panic recovery keeps the forum from crashing on unexpected panics.
	// NOTE: if a panic happens after part of the response was already written,
	// the status/body may already be partially committed by net/http.
	var handler http.Handler = mux
	handler = withNotFound(mux)
	handler = handlers.SecurityHeaders(handler)
	handler = recoverPanic(handler)

	return &Server{
		handler: handler,
	}
}

// Run starts HTTP (redirect) and HTTPS servers.
// AUDIT: enforces HTTPS by redirecting all HTTP traffic.
func (s *Server) Run() error {
	certFile := os.Getenv("TLS_CERT_FILE")
	keyFile := os.Getenv("TLS_KEY_FILE")

	if certFile == "" {
		certFile = "./certs/cert.pem"
	}
	if keyFile == "" {
		keyFile = "./certs/key.pem"
	}

	httpsAddr := ":8443"
	httpAddr := ":8080"

	go func() {
		log.Println("HTTP redirect server running at http://localhost" + httpAddr)

		err := http.ListenAndServe(httpAddr, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			host := r.Host

			switch {
			case strings.HasPrefix(host, "localhost:"):
				host = "localhost:8443"
			case strings.HasPrefix(host, "127.0.0.1:"):
				host = "127.0.0.1:8443"
			case host == "localhost":
				host = "localhost:8443"
			case host == "127.0.0.1":
				host = "127.0.0.1:8443"
			}

			target := "https://" + host + r.URL.RequestURI()
			http.Redirect(w, r, target, http.StatusMovedPermanently)
		}))
		if err != nil {
			log.Println("HTTP redirect server stopped:", err)
		}
	}()

	srv := &http.Server{
		Addr:      httpsAddr,
		Handler:   s.handler,
		TLSConfig: handlers.SecureTLSConfig(),
	}

	log.Println("HTTPS server running at https://localhost" + httpsAddr)
	return srv.ListenAndServeTLS(certFile, keyFile)
}