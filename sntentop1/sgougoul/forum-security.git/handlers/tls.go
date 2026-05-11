package handlers

import "crypto/tls"

// SecureTLSConfig returns a reasonable TLS configuration for local HTTPS use.
// AUDIT:
// - MinVersion blocks obsolete TLS versions
// - CurvePreferences prefers modern elliptic curves
// - PreferServerCipherSuites is ignored in TLS 1.3 but harmless here
func SecureTLSConfig() *tls.Config {
	return &tls.Config{
		MinVersion: tls.VersionTLS12,

		// AUDIT: modern curve preference for TLS handshakes.
		CurvePreferences: []tls.CurveID{
			tls.X25519,
			tls.CurveP256,
		},

		// AUDIT: kept for compatibility with TLS 1.2 behavior.
		PreferServerCipherSuites: true,
	}
}