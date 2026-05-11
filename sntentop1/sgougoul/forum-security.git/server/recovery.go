package server

import (
	"log"
	"net/http"

	"forum/handlers"
)

// recoverPanic catches unexpected panics and returns a styled 500 page.
// AUDIT: prevents request crashes from exposing raw panic behavior to users.
func recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				log.Printf("panic recovered: %v", rec)
				handlers.RenderError(w, r, http.StatusInternalServerError, "Internal Server Error")
			}
		}()

		next.ServeHTTP(w, r)
	})
}