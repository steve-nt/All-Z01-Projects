package middleware

import (
	"log"
	"net/http"
	"social-network/backend/utils"
)

// HandlerFunc aliases http.HandlerFunc for easier middleware composition
type HandlerFunc = http.HandlerFunc

// LoggingMiddleware logs all HTTP requests
func LoggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Create a wrapper to capture status code
		lrw := &loggingResponseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		// Log the request
		log.Printf("%s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)

		// Call the next handler
		next(lrw, r)

		// Log the response status
		log.Printf("%s %s - Status: %d", r.Method, r.URL.Path, lrw.statusCode)
	}
}

// loggingResponseWriter wraps http.ResponseWriter to capture status code
type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

// RequireAuth wraps a handler to require authentication
// If user is not authenticated, redirects to login page
// Use this for HTML page routes
func RequireAuth(next HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Check if user is authenticated
		isAuth, userID, _ := utils.CheckAuth(r)
		if !isAuth || userID == 0 {
			// Redirect to login if not authenticated
			http.Redirect(w, r, "/login?error=Please log in to access this page", http.StatusSeeOther)
			return
		}

		// User is authenticated, call the next handler
		next(w, r)
	}
}

// RequireAuthJSON wraps a handler to require authentication for API endpoints
// Returns JSON error instead of redirect
// Use this for API routes that return JSON
func RequireAuthJSON(next HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Check if user is authenticated
		isAuth, userID, _ := utils.CheckAuth(r)
		if !isAuth || userID == 0 {
			// Return JSON error for API endpoints
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"error": "Unauthorized", "message": "Please log in to access this resource"}`))
			return
		}

		// User is authenticated, call the next handler
		next(w, r)
	}
}

// ErrorHandlingMiddleware catches panics and serves 500 error page
func ErrorHandlingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic occurred: %v", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()
		next(w, r)
	}
}

// WrapHandler is a convenience function that applies common middleware
// (logging + error handling) to a handler
func WrapHandler(handler HandlerFunc) http.HandlerFunc {
	return LoggingMiddleware(ErrorHandlingMiddleware(handler))
}

// CombinedMiddleware applies multiple middlewares in order
// You can chain middlewares: Logging -> ErrorHandling -> Auth
func CombinedMiddleware(middlewares ...func(http.HandlerFunc) http.HandlerFunc) func(HandlerFunc) http.HandlerFunc {
	return func(next HandlerFunc) http.HandlerFunc {
		handler := http.HandlerFunc(next)
		// Apply middlewares in reverse order (last middleware wraps first)
		for i := len(middlewares) - 1; i >= 0; i-- {
			handler = middlewares[i](handler)
		}
		return handler
	}
}
