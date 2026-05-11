package handlers

import (
	"net"
	"net/http"
	"sync"
	"time"
)

// visitor stores per-IP rate limiting state.
type visitor struct {
	tokens     float64
	lastSeen   time.Time
	lastRefill time.Time
}

// rateLimiter provides a simple in-memory token bucket limiter.
// AUDIT: this protects authentication and write actions from burst abuse.
type rateLimiter struct {
	mu              sync.Mutex
	visitors        map[string]*visitor
	maxTokens       float64
	refillPerSecond float64
}

// newRateLimiter creates a token bucket limiter.
// maxTokens controls burst size.
// refillPerSecond controls sustained rate.
func newRateLimiter(maxTokens, refillPerSecond float64) *rateLimiter {
	rl := &rateLimiter{
		visitors:        make(map[string]*visitor),
		maxTokens:       maxTokens,
		refillPerSecond: refillPerSecond,
	}

	// AUDIT: cleanup loop prevents stale visitor entries from accumulating forever.
	go rl.cleanupLoop()

	return rl
}

// allow reports whether the request from this IP should be allowed.
func (rl *rateLimiter) allow(ip string) bool {
	now := time.Now()

	rl.mu.Lock()
	defer rl.mu.Unlock()

	v, ok := rl.visitors[ip]
	if !ok {
		rl.visitors[ip] = &visitor{
			tokens:     rl.maxTokens - 1,
			lastSeen:   now,
			lastRefill: now,
		}
		return true
	}

	// Refill tokens based on elapsed time.
	elapsed := now.Sub(v.lastRefill).Seconds()
	v.tokens += elapsed * rl.refillPerSecond
	if v.tokens > rl.maxTokens {
		v.tokens = rl.maxTokens
	}

	v.lastSeen = now
	v.lastRefill = now

	if v.tokens < 1 {
		return false
	}

	v.tokens -= 1
	return true
}

// cleanupLoop periodically removes inactive IP entries.
func (rl *rateLimiter) cleanupLoop() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		cutoff := time.Now().Add(-15 * time.Minute)

		rl.mu.Lock()
		for ip, v := range rl.visitors {
			if v.lastSeen.Before(cutoff) {
				delete(rl.visitors, ip)
			}
		}
		rl.mu.Unlock()
	}
}

// clientIP extracts the remote IP address from the request.
// For local dev and direct deployments, RemoteAddr is enough.
func clientIP(r *http.Request) string {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}

// RateLimit wraps a handler with IP-based rate limiting.
func RateLimit(next http.HandlerFunc, maxTokens, refillPerSecond float64) http.HandlerFunc {
	rl := newRateLimiter(maxTokens, refillPerSecond)

	return func(w http.ResponseWriter, r *http.Request) {
		ip := clientIP(r)

		if !rl.allow(ip) {
			// AUDIT: 429 is the correct status for rate limiting.
			RenderError(w, r, http.StatusTooManyRequests, "Too many requests. Please slow down and try again.")
			return
		}

		next(w, r)
	}
}