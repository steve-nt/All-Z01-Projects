package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"syscall"

	"ascii/ascii"
	"ascii/web"
)

// Custom 404 handler
func notFoundHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, `<html>
        <head>
            <link rel="stylesheet" type="text/css" href="/styles.css">
        </head>
        <body>
            <div class="error-container">
				<h1>404 - Page Not Found</h1>
				<h2>🤔</h2>
				<p>The page you are looking for does not exist.</p>
			<a href="/">Take me back to ASCII Art</a>
			 </div>
        </body>
    </html>`)
}

// Custom 400 handler for banner or text errors
func badRequestHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusBadRequest)
	fmt.Fprintf(w, `<html>
        <head>
            <link rel="stylesheet" type="text/css" href="/styles.css">
        </head>
        <body>
            <div class="error-container">
                <h1>400 - Bad Request</h1>
				<h2>👾</h2>
                <p>Your request is invalid or malformed.</p>
                <a href="/">Take me back to ASCII Art</a>
            </div>
        </body>
    </html>`)
}

// Custom 500 handler
func internalServerErrorHandler(w http.ResponseWriter, r *http.Request) {
	// Set the content type to HTML
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	// Set the HTTP status to 500 (Internal Server Error)
	w.WriteHeader(http.StatusInternalServerError)
	// Respond with a meaningful error message
	fmt.Fprintf(w, `<html>
        <head>
            <link rel="stylesheet" type="text/css" href="/styles.css">
        </head>
        <body>
            <div class="error-container">
                <h1>500 - Internal Server Error</h1>
				<h2>🤷</h2>
			<p>An unexpected error occurred on the server. Please try again later.</p>
			<a href="/">Take me back to ASCII Art</a>
            </div>
        </body>
    </html>`)
}

func main() {
	err := ascii.HandleArgs(os.Args)
	if err != nil {
		fmt.Println(err)
	}

	// Handle CSS file requests
	http.HandleFunc("/styles.css", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/css; charset=utf-8")
		http.ServeFile(w, r, "templates/styles.css")
	})

	// Routes for testing 400 and 500 errors
	http.HandleFunc("/bad-request", badRequestHandler)
	http.HandleFunc("/internal-error", internalServerErrorHandler)

	// General route handler
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		text := r.FormValue("text")
		// Validate that the input contains only ASCII characters (0-127)
		asciiPattern := regexp.MustCompile(`^[\x00-\x7F]*$`)

		// Check if the text contains non-ASCII characters
		if !asciiPattern.MatchString(text) {
			// Return a 400 Bad Request if the path contains unsupported characters
			badRequestHandler(w, r)
		} else if r.URL.Path != "/" {
			notFoundHandler(w, r) // Handle undefined paths
		} else {
			web.HandleHTTP(w, r) // Handle the root path
		}
	})

	// Create HTTP server
	server := &http.Server{
		Addr: ":8080",
	}

	// Run the server in a goroutine
	go func() {
		fmt.Println("🚀 Server running at http://localhost:8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("❌ Server error: %v\n", err)
			os.Exit(1)
		}
	}()

	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	fmt.Println("\n⚙️  Shutting down server...")
	if err := server.Close(); err != nil {
		fmt.Printf("❌ Error during server shutdown: %v\n", err)
	} else {
		fmt.Println("✅ Server gracefully stopped.")
	}
}
