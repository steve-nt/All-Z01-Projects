package websocket

import (
	"log"
	"net/http"
)

// SetupWebSocketRoutes registers WebSocket endpoints
func SetupWebSocketRoutes(hub *Hub) {
	// WebSocket endpoint for real-time communication
	// Note: We don't use middleware here because WebSocket upgrades must happen
	// before any response is written. Middleware can interfere with the upgrade process.
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		// Simple logging without middleware
		log.Printf("WebSocket connection attempt from %s", r.RemoteAddr)
		
		// Handle WebSocket upgrade
		ServeWS(hub, w, r)
	})
}

