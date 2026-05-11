package websocket

import (
	"log"
	"sync"
)

// Global hub instance (set by main.go)
var globalHub *Hub

// SetGlobalHub sets the global hub instance
func SetGlobalHub(hub *Hub) {
	globalHub = hub
}

// GetGlobalHub returns the global hub instance
func GetGlobalHub() *Hub {
	return globalHub
}

// Hub maintains the set of active clients and broadcasts messages to clients
type Hub struct {
	// Registered clients mapped by user ID
	clients map[int]*Client

	// Mutex to protect clients map
	mu sync.RWMutex

	// Inbound messages from the clients
	broadcast chan *Message

	// Register requests from the clients
	register chan *Client

	// Unregister requests from clients
	unregister chan *Client
}

// Message represents a WebSocket message
type Message struct {
	Type      string      `json:"type"`       // message type: "private_message", "group_message", "notification", etc.
	UserID    int         `json:"user_id"`    // sender user ID (for messages)
	Recipient int         `json:"recipient"`  // recipient user ID (for private messages)
	GroupID   int         `json:"group_id"`   // group ID (for group messages)
	Data      interface{} `json:"data"`       // message payload
}

// NewHub creates a new Hub instance
func NewHub() *Hub {
	return &Hub{
		clients:    make(map[int]*Client),
		broadcast:  make(chan *Message, 256),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

// Run starts the hub's message loop
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			if oldClient, exists := h.clients[client.userID]; exists {
				log.Printf("Closing existing connection for user %d", client.userID)
				close(oldClient.send)
			}
			h.clients[client.userID] = client
			// Tell the new client who else is online; tell everyone else the new client is online
			for otherID, otherClient := range h.clients {
				if otherID == client.userID {
					continue
				}
				presenceNew := &Message{Type: "presence", Data: map[string]interface{}{"user_id": otherID, "online": true}}
				select {
				case client.send <- presenceNew:
				default:
				}
				presenceMe := &Message{Type: "presence", Data: map[string]interface{}{"user_id": client.userID, "online": true}}
				select {
				case otherClient.send <- presenceMe:
				default:
					close(otherClient.send)
					delete(h.clients, otherID)
				}
			}
			h.mu.Unlock()
			log.Printf("Client registered: user_id=%d, total clients=%d", client.userID, len(h.clients))

		case client := <-h.unregister:
			h.mu.Lock()
			userID := client.userID
			if _, ok := h.clients[userID]; ok {
				delete(h.clients, userID)
				close(client.send)
				// Tell remaining clients this user went offline
				offlineMsg := &Message{Type: "presence", Data: map[string]interface{}{"user_id": userID, "online": false}}
				for _, c := range h.clients {
					select {
					case c.send <- offlineMsg:
					default:
					}
				}
				log.Printf("Client unregistered: user_id=%d, total clients=%d", userID, len(h.clients))
			}
			h.mu.Unlock()

		case message := <-h.broadcast:
			h.mu.RLock()
			switch message.Type {
			case "private_message":
				// Send to specific recipient
				if client, ok := h.clients[message.Recipient]; ok {
					select {
					case client.send <- message:
					default:
						close(client.send)
						delete(h.clients, message.Recipient)
					}
				}
			case "group_message":
				// Broadcast to all group members (we'll need to get group members from DB)
				// For now, we'll handle this in the handler that calls BroadcastToGroup
				for userID, client := range h.clients {
					// Skip sender
					if userID == message.UserID {
						continue
					}
					// TODO: Check if user is member of the group (will be handled in handler)
					select {
					case client.send <- message:
					default:
						close(client.send)
						delete(h.clients, userID)
					}
				}
			case "notification":
				// Send notification to specific user
				if client, ok := h.clients[message.Recipient]; ok {
					select {
					case client.send <- message:
					default:
						close(client.send)
						delete(h.clients, message.Recipient)
					}
				}
			default:
				// Broadcast to all clients
				for userID, client := range h.clients {
					select {
					case client.send <- message:
					default:
						close(client.send)
						delete(h.clients, userID)
					}
				}
			}
			h.mu.RUnlock()
		}
	}
}

// GetClient returns a client for a given user ID
func (h *Hub) GetClient(userID int) (*Client, bool) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	client, ok := h.clients[userID]
	return client, ok
}

// BroadcastToUser sends a message to a specific user
func (h *Hub) BroadcastToUser(userID int, message *Message) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	if client, ok := h.clients[userID]; ok {
		select {
		case client.send <- message:
		default:
			close(client.send)
			delete(h.clients, userID)
		}
	}
}

// BroadcastToGroup sends a message to all members of a group (except sender)
// Note: This is a placeholder - the actual group membership check will be done in handlers
func (h *Hub) BroadcastToGroup(groupID int, senderID int, message *Message) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	// This will be enhanced in handlers to check actual group membership
	for userID, client := range h.clients {
		if userID == senderID {
			continue
		}
		select {
		case client.send <- message:
		default:
			close(client.send)
			delete(h.clients, userID)
		}
	}
}

// GetConnectedUsers returns a list of currently connected user IDs
func (h *Hub) GetConnectedUsers() []int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	users := make([]int, 0, len(h.clients))
	for userID := range h.clients {
		users = append(users, userID)
	}
	return users
}

