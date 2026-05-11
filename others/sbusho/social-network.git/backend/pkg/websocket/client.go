package websocket

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"social-network/backend/utils"
)

const (
	// Time allowed to write a message to the peer
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer
	pongWait = 60 * time.Second

	// Send pings to peer with this period (must be less than pongWait)
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer
	maxMessageSize = 512 * 1024 // 512 KB
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// Allow connections from any origin (for development)
		// In production, you should check the origin
		return true
	},
}

// Client is a middleman between the websocket connection and the hub
type Client struct {
	hub *Hub

	// The websocket connection
	conn *websocket.Conn

	// Buffered channel of outbound messages
	send chan *Message

	// User ID of this client
	userID int
}

// readPump pumps messages from the websocket connection to the hub
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, messageBytes, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}

		// Parse incoming message
		var msg Message
		if err := json.Unmarshal(messageBytes, &msg); err != nil {
			log.Printf("Error parsing message: %v", err)
			continue
		}

		// Handle different message types here if needed
		log.Printf("Received message from user %d: type=%s", c.userID, msg.Type)

		// For testing: Echo back test messages
		if msg.Type == "test" {
			echoMsg := &Message{
				Type:   "echo",
				UserID: c.userID,
				Data: map[string]interface{}{
					"original": msg.Data,
					"message":  "Message received and echoed back!",
					"timestamp": time.Now().Format(time.RFC3339),
				},
			}
			select {
			case c.send <- echoMsg:
			default:
				log.Printf("Failed to send echo message to user %d", c.userID)
			}
		}
	}
}

// writePump pumps messages from the hub to the websocket connection
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}

			jsonData, err := json.Marshal(message)
			if err != nil {
				log.Printf("Error marshaling message: %v", err)
				w.Close()
				continue
			}

			w.Write(jsonData)

			// Add queued messages to the current websocket message
			n := len(c.send)
			for i := 0; i < n; i++ {
				queuedMsg := <-c.send
				w.Write([]byte{'\n'})
				jsonData, err := json.Marshal(queuedMsg)
				if err != nil {
					log.Printf("Error marshaling queued message: %v", err)
					continue
				}
				w.Write(jsonData)
			}

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// ServeWS handles websocket requests from clients
func ServeWS(hub *Hub, w http.ResponseWriter, r *http.Request) {
	// Authenticate the connection using existing session system
	cookie, err := r.Cookie("session")
	if err != nil {
		log.Printf("WebSocket connection rejected: no session cookie (error: %v)", err)
		// For WebSocket, we need to send a proper close frame
		// But since upgrade hasn't happened, we can't. So we'll let the upgrade fail.
		http.Error(w, "Unauthorized: No session cookie", http.StatusUnauthorized)
		return
	}

	log.Printf("WebSocket connection attempt: cookie present, validating session...")

	if !utils.IsValidSession(cookie.Value) {
		cookiePreview := cookie.Value
		if len(cookiePreview) > 10 {
			cookiePreview = cookiePreview[:10] + "..."
		}
		log.Printf("WebSocket connection rejected: invalid session (cookie value: %s)", cookiePreview)
		http.Error(w, "Unauthorized: Invalid session", http.StatusUnauthorized)
		return
	}

	userID := utils.GetUserIDFromSession(cookie.Value)
	if userID == 0 {
		log.Printf("WebSocket connection rejected: invalid user ID")
		http.Error(w, "Unauthorized: Invalid user ID", http.StatusUnauthorized)
		return
	}

	log.Printf("WebSocket authentication successful for user %d, upgrading connection...", userID)

	// Upgrade connection to WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}

	// Create new client
	client := &Client{
		hub:    hub,
		conn:   conn,
		send:   make(chan *Message, 256),
		userID: userID,
	}

	// Register client with hub
	client.hub.register <- client

	log.Printf("WebSocket connection established for user %d", userID)

	// Start client goroutines
	go client.writePump()
	go client.readPump()
}

