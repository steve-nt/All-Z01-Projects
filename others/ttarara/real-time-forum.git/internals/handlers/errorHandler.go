package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"realtimeforum/internals/utils"
)

// NotFoundHandler handles 404 errors
func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	// Check if it's an API request
	if utils.IsAPIRequest(r.URL.Path) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"error": true, "message": "Not Found", "status": 404}`))
		return
	}
	// For SPA, serve index.html
	http.ServeFile(w, r, "frontend/index.html")
}

// InternalServerErrorHandler handles 500 errors
func InternalServerErrorHandler(w http.ResponseWriter, r *http.Request) {
	// Check if it's an API request
	if utils.IsAPIRequest(r.URL.Path) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": true, "message": "Internal Server Error", "status": 500}`))
		return
	}
	// For SPA, serve index.html
	http.ServeFile(w, r, "frontend/index.html")
}

// BadRequestHandler handles 400 errors
func BadRequestHandler(w http.ResponseWriter, r *http.Request) {
	// Check if it's an API request
	if utils.IsAPIRequest(r.URL.Path) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": true, "message": "Bad Request", "status": 400}`))
		return
	}
	// For SPA, serve index.html
	http.ServeFile(w, r, "frontend/index.html")
}

// UnauthorizedHandler handles 401 errors
func UnauthorizedHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusUnauthorized)
	// Redirect to login for unauthorized access
	http.Redirect(w, r, "/login?error=unauthorized", http.StatusSeeOther)
}

// ForbiddenHandler handles 403 errors
func ForbiddenHandler(w http.ResponseWriter, r *http.Request) {
	// Check if it's an API request
	if utils.IsAPIRequest(r.URL.Path) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(`{"error": true, "message": "Access Forbidden", "status": 403}`))
		return
	}
	// For SPA, serve index.html
	http.ServeFile(w, r, "frontend/index.html")
}

// CustomErrorHandler is a middleware to catch errors and serve appropriate pages
func CustomErrorHandler(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic occurred on %s %s: %v", r.Method, r.URL.Path, err)
				InternalServerErrorHandler(w, r)
			}
		}()
		next(w, r)
	}
}

// ErrorResponse sends a JSON error response for API endpoints
func ErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	// Create and send JSON response
	response := map[string]interface{}{
		"error":   true,
		"message": message,
		"status":  statusCode,
	}

	json.NewEncoder(w).Encode(response)
}

// HandleHTTPError handles different HTTP errors based on status code
func HandleHTTPError(w http.ResponseWriter, r *http.Request, statusCode int, message string) {
	// Check if this is an API request (JSON response expected)
	if r.Header.Get("Content-Type") == "application/json" ||
		r.Header.Get("Accept") == "application/json" ||
		r.URL.Path[:4] == "/api" {
		ErrorResponse(w, statusCode, message)
		return
	}

	// Handle HTML error pages
	switch statusCode {
	case http.StatusNotFound:
		NotFoundHandler(w, r)
	case http.StatusBadRequest:
		BadRequestHandler(w, r)
	case http.StatusUnauthorized:
		UnauthorizedHandler(w, r)
	case http.StatusForbidden:
		ForbiddenHandler(w, r)
	case http.StatusInternalServerError:
		fallthrough
	default:
		InternalServerErrorHandler(w, r)
	}
}

// NotFoundMiddleware catches all unmatched routes
func NotFoundMiddleware() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		NotFoundHandler(w, r)
	}
}
