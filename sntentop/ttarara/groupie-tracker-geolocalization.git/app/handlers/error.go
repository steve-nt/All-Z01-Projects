package handlers

import (
	"html/template"
	"log"
	"net/http"
	"strconv"
)

type ErrorData struct {
	StatusCode int
	StatusText string
	Message    string
}

// renderError handles all error page rendering
func renderError(w http.ResponseWriter, status int, r *http.Request, message ...string) {
	// Determine template path
	templatePath := "web/templates/errors/" + strconv.Itoa(status) + ".html"

	// Prepare error data
	data := ErrorData{
		StatusCode: status,
		StatusText: http.StatusText(status),
	}

	if len(message) > 0 {
		data.Message = message[0]
	}

	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		log.Printf("Failed to parse error template: %v", err)
		http.Error(w, data.StatusText, status)
		return
	}

	w.WriteHeader(status)
	if err := tmpl.Execute(w, data); err != nil {
		log.Printf("Error rendering template: %v", err)
		http.Error(w, data.StatusText, status)
	}
}

// LoggingMiddleware logs each incoming HTTP request method and path
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Received request: %s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

// RecoveryMiddleware recovers from panics and responds with a 500 error page
func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Recovered from panic: %v", err)
				renderError(w, http.StatusInternalServerError, r, "Internal Server Error")
			}
		}()
		next.ServeHTTP(w, r)
	})
}
