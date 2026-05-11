package server

import (
	"net/http"

	"forum/handlers"
)

// withNotFound returns a custom 404 page.
// AUDIT: ensures proper HTTP status for unknown routes.
func withNotFound(mux *http.ServeMux) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			mux.ServeHTTP(w, r)
			return
		}

		_, pattern := mux.Handler(r)
		if pattern == "" || pattern == "/" {
			handlers.RenderError(w, r, http.StatusNotFound, "Page not found.")
			return
		}

		mux.ServeHTTP(w, r)
	})
}