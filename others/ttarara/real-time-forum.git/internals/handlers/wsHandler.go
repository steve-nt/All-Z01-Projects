package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"realtimeforum/internals/utils"

	"github.com/gorilla/websocket"
)


var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Client represents a WebSocket connection with user information
type Client struct {
	conn     *websocket.Conn
	userID   int
	username string
}

var (
	wsClients   = make(map[*websocket.Conn]*Client)
	wsBroadcast = make(chan []byte)
)

// IsUserOnline checks if a user is currently connected via WebSocket
func IsUserOnline(userID int) bool {
	for _, client := range wsClients {
		if client.userID == userID {
			return true
		}
	}
	return false
}

// SendNotificationToUser sends a notification to a specific user via WebSocket
func SendNotificationToUser(userID int, notification map[string]interface{}) {
	notificationJSON, err := json.Marshal(notification)
	if err != nil {
		log.Printf("Error marshaling notification: %v\n", err)
		return
	}

	// Send to all connections for this user
	sent := false
	for conn, client := range wsClients {
		if client.userID == userID {
			if err := conn.WriteMessage(websocket.TextMessage, notificationJSON); err != nil {
				log.Printf("WebSocket write error to user %d: %v\n", userID, err)
				conn.Close()
				delete(wsClients, conn)
			} else {
				sent = true
				log.Printf("Notification sent to user %d (username: %s)\n", userID, client.username)
			}
		}
	}
	if !sent && userID > 0 {
		log.Printf("Warning: User %d not found in connected clients for notification\n", userID)
	}
}


func init() {
	go wsHandleMessages()
}


func WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	// Αναβάθμιση της σύνδεσης σε WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v\n", err)
		return
	}
	defer conn.Close()

	// Get user info from session cookie
	var userID int
	var username string
	if cookie, err := r.Cookie("session"); err == nil && utils.IsValidSession(cookie.Value) {
		userID = utils.GetUserIDFromSession(cookie.Value)
		username = utils.GetUsernameFromSession(cookie.Value)
	}

	// Create client
	client := &Client{
		conn:     conn,
		userID:   userID,
		username: username,
	}


	wsClients[conn] = client
	log.Printf("New WebSocket client connected: userID=%d, username=%s, total_clients=%d\n", userID, username, len(wsClients))
	
	// Broadcast user status update when user connects
	if userID > 0 {
		broadcastUserStatusUpdate()
	}


	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Printf("WebSocket read error: %v\n", err)
			delete(wsClients, conn)
			// Broadcast user status update when user disconnects
			if userID > 0 {
				broadcastUserStatusUpdate()
			}
			break
		}

		wsBroadcast <- msg
	}
}

func broadcastUserStatusUpdate() {
	// Get all unique user IDs that are online
	onlineUserIDs := make(map[int]bool)
	for _, client := range wsClients {
		if client.userID > 0 {
			onlineUserIDs[client.userID] = true
		}
	}
	
	// Create status update message
	statusUpdate := map[string]interface{}{
		"type":  "user_status_update",
		"users":  onlineUserIDs,
	}
	
	if jsonStatus, err := json.Marshal(statusUpdate); err == nil {
		// Send to all connected clients
		for conn := range wsClients {
			if err := conn.WriteMessage(websocket.TextMessage, jsonStatus); err != nil {
				log.Printf("Error sending user status update: %v\n", err)
				conn.Close()
				delete(wsClients, conn)
			}
		}
	}
}


func wsHandleMessages() {
	for {
		msg := <-wsBroadcast

		// Try to parse as JSON to check message type
		var msgData map[string]interface{}
		if err := json.Unmarshal(msg, &msgData); err == nil {
			// Check if it's a private message
			if msgType, ok := msgData["type"].(string); ok && msgType == "private_message" {
				// Get receiver and sender IDs (JSON numbers come as float64)
				var receiverIDInt, senderIDInt int
				
				// Handle receiver_id (can be float64 from JSON or int)
				if receiverID, ok := msgData["receiver_id"].(float64); ok {
					receiverIDInt = int(receiverID)
				} else if receiverID, ok := msgData["receiver_id"].(int); ok {
					receiverIDInt = receiverID
				}
				
				// Handle sender_id (can be float64 from JSON or int)
				if senderID, ok := msgData["sender_id"].(float64); ok {
					senderIDInt = int(senderID)
				} else if senderID, ok := msgData["sender_id"].(int); ok {
					senderIDInt = senderID
				}
				
				log.Printf("Sending private message: sender_id=%d, receiver_id=%d, total_clients=%d\n", 
					senderIDInt, receiverIDInt, len(wsClients))
				
				// Send to receiver
				if receiverIDInt > 0 {
					sentToReceiver := false
					for conn, client := range wsClients {
						if client.userID == receiverIDInt {
							if err := conn.WriteMessage(websocket.TextMessage, msg); err != nil {
								log.Printf("WebSocket write error to receiver %d: %v\n", receiverIDInt, err)
								conn.Close()
								delete(wsClients, conn)
							} else {
								sentToReceiver = true
								log.Printf("Message sent to receiver %d (username: %s)\n", receiverIDInt, client.username)
							}
						}
					}
					if !sentToReceiver {
						log.Printf("Warning: Receiver %d not found in connected clients\n", receiverIDInt)
					}
				}
				
				// Also send to sender for confirmation
				if senderIDInt > 0 {
					sentToSender := false
					for conn, client := range wsClients {
						if client.userID == senderIDInt {
							if err := conn.WriteMessage(websocket.TextMessage, msg); err != nil {
								log.Printf("WebSocket write error to sender %d: %v\n", senderIDInt, err)
								conn.Close()
								delete(wsClients, conn)
							} else {
								sentToSender = true
								log.Printf("Message sent to sender %d (username: %s)\n", senderIDInt, client.username)
							}
						}
					}
					if !sentToSender {
						log.Printf("Warning: Sender %d not found in connected clients\n", senderIDInt)
					}
				}
				
				continue
			}
			
			// Handle typing indicators
			if msgType, ok := msgData["type"].(string); ok && (msgType == "typing_start" || msgType == "typing_stop") {
				// Get receiver and sender IDs
				var receiverIDInt, senderIDInt int
				
				if receiverID, ok := msgData["receiver_id"].(float64); ok {
					receiverIDInt = int(receiverID)
				} else if receiverID, ok := msgData["receiver_id"].(int); ok {
					receiverIDInt = receiverID
				}
				
				if senderID, ok := msgData["sender_id"].(float64); ok {
					senderIDInt = int(senderID)
				} else if senderID, ok := msgData["sender_id"].(int); ok {
					senderIDInt = senderID
				}
				
				// Only send to receiver (not sender)
				if receiverIDInt > 0 {
					log.Printf("Typing indicator: sender_id=%d, receiver_id=%d, type=%s\n", senderIDInt, receiverIDInt, msgType)
					for conn, client := range wsClients {
						if client.userID == receiverIDInt {
							if err := conn.WriteMessage(websocket.TextMessage, msg); err != nil {
								log.Printf("WebSocket write error to receiver %d: %v\n", receiverIDInt, err)
								conn.Close()
								delete(wsClients, conn)
							}
						}
					}
				}
				continue
			}
		}

		// Default: broadcast to all clients (for posts, comments, etc.)
		for conn := range wsClients {
			if err := conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				log.Printf("WebSocket write error: %v\n", err)
				conn.Close()
				delete(wsClients, conn)
			}
		}
	}
}

