package handlers

import "net/http"

// SecurityHeaders adds a baseline set of defensive HTTP headers.
// AUDIT:
// - nosniff reduces MIME confusion attacks
// - DENY prevents clickjacking through framing
// - referrer policy limits referrer leakage
// - CSP restricts where active content may load from
func SecurityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Prevent MIME-type sniffing.
		w.Header().Set("X-Content-Type-Options", "nosniff")

		// Prevent the site from being embedded in frames/iframes.
		w.Header().Set("X-Frame-Options", "DENY")

		// Limit how much referrer information browsers send.
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")

		// Minimal CSP suitable for this server-rendered forum.
		// AUDIT: resources are restricted to the same origin.
		w.Header().Set(
			"Content-Security-Policy",
			"default-src 'self'; style-src 'self' 'unsafe-inline'; img-src 'self' data:; script-src 'self'; object-src 'none'; base-uri 'self'; frame-ancestors 'none'; form-action 'self'",
		)

		next.ServeHTTP(w, r)
	})
}