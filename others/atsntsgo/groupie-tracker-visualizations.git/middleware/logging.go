package middleware

import (
	"log"
	"net/http"
	"time"
)

func LoggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Create a custom response writer to capture the status code
		crw := &customResponseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		// Log the incoming request
		startTime := time.Now()
		log.Printf("Request Started - Method: %s, Path: %s, Client IP: %s",
			r.Method, r.URL.Path, r.RemoteAddr)

		// Call the next handler
		next.ServeHTTP(crw, r)

		// Log the completion with status code and duration
		duration := time.Since(startTime)
		log.Printf("Request Completed - Status: %d, Duration: %v, Path: %s",
			crw.statusCode, duration, r.URL.Path)
	}
}

// customResponseWriter wraps http.ResponseWriter to capture the status code
type customResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (crw *customResponseWriter) WriteHeader(code int) {
	crw.statusCode = code
	crw.ResponseWriter.WriteHeader(code)
}
