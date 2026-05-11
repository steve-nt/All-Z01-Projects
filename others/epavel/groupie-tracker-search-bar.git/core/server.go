package core

import (
	"context"
	"fmt"
	"groupie-tracker/bin"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
)

// SetupServer sets up the server with the given handler and custom middleware for advanced handling
func SetupServer(handler http.Handler) *http.Server {
	port := FindOpenPort(defaultPort)
	addr := fmt.Sprintf(":%d", port)
	return &http.Server{
		Addr:    addr,
		Handler: bin.Middleware(handler),
	}
}

// FindOpenPort finds an open port starting from the given port
func FindOpenPort(startPort int) int {
	port := startPort
	for {
		addr := fmt.Sprintf(":%d", port)
		listener, err := net.Listen("tcp", addr)
		if err == nil {
			listener.Close()
			return port
		}
		port++
	}
}

// StartServer starts the server and listens for interrupt or terminate signals
func StartServer(server *http.Server) {
	// Channel to listen for interrupt or terminate signals
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	// Run server in a goroutine
	go func() {
		fmt.Printf("Server is running on http://localhost%s\n", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Could not listen on %s: %v\n", server.Addr, err)
		}
	}()

	// Block until we receive a signal
	<-stop

	// Create a deadline to wait for server shutdown
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	// Attempt graceful shutdown
	fmt.Println("\nShutting down server...")
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	fmt.Println("Server exiting")
}
