package messages

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"social-network/backend/pkg/db/sqlite"
	"social-network/backend/pkg/websocket"
	"social-network/backend/utils"
)

type MessageResponse struct {
	ID          int    `json:"id"`
	SenderID    int    `json:"sender_id"`
	SenderName  string `json:"sender_name"`
	RecipientID *int   `json:"recipient_id,omitempty"`
	GroupID     *int   `json:"group_id,omitempty"`
	Content     string `json:"content"`
	CreatedAt   string `json:"created_at"`
	ReadAt      *string `json:"read_at,omitempty"`
}

// canSendMessage checks if sender can send a message to recipient
// Returns true if:
// - At least one user follows the other (accepted status)
// - OR recipient has a public profile
func canSendMessage(db *sql.DB, senderID int, recipientID int) (bool, error) {
	// Check if recipient has public profile
	var isPublic bool
	err := db.QueryRow(`SELECT is_public FROM Users WHERE user_id = ?`, recipientID).Scan(&isPublic)
	if err != nil {
		return false, err
	}
	if isPublic {
		return true, nil
	}

	// Check if at least one follows the other (accepted status)
	var canMessage bool
	err = db.QueryRow(`
		SELECT EXISTS(
			SELECT 1 FROM Followers
			WHERE (
				(follower_id = ? AND following_id = ? AND status = 'accepted') OR
				(follower_id = ? AND following_id = ? AND status = 'accepted')
			)
		)
	`, senderID, recipientID, recipientID, senderID).Scan(&canMessage)
	if err != nil {
		return false, err
	}

	return canMessage, nil
}

// canReceiveMessage checks if recipient can receive message from sender via WebSocket
// Returns true if:
// - Recipient is following sender (accepted status)
// - OR recipient has a public profile
func canReceiveMessage(db *sql.DB, senderID int, recipientID int) (bool, error) {
	// Check if recipient has public profile
	var isPublic bool
	err := db.QueryRow(`SELECT is_public FROM Users WHERE user_id = ?`, recipientID).Scan(&isPublic)
	if err != nil {
		return false, err
	}
	if isPublic {
		return true, nil
	}

	// Check if recipient is following sender (accepted status)
	var canReceive bool
	err = db.QueryRow(`
		SELECT EXISTS(
			SELECT 1 FROM Followers
			WHERE follower_id = ? AND following_id = ? AND status = 'accepted'
		)
	`, recipientID, senderID).Scan(&canReceive)
	if err != nil {
		return false, err
	}

	return canReceive, nil
}

// SendMessageHandler handles POST /api/messages/send
// Sends a private message to another user
func SendMessageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSONError(w, http.StatusMethodNotAllowed, "method_not_allowed", "Method not allowed")
		return
	}

	// Get authenticated user
	isAuth, senderID, _ := utils.CheckAuth(r)
	if !isAuth || senderID == 0 {
		writeJSONError(w, http.StatusUnauthorized, "unauthorized", "Please log in to send messages")
		return
	}

	var payload struct {
		RecipientID int    `json:"recipient_id"`
		Content     string `json:"content"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid_json", "Invalid JSON body")
		return
	}

	// Validate input
	if payload.RecipientID <= 0 {
		writeJSONError(w, http.StatusBadRequest, "invalid_request", "recipient_id is required")
		return
	}
	if payload.RecipientID == senderID {
		writeJSONError(w, http.StatusBadRequest, "invalid_request", "Cannot send message to yourself")
		return
	}

	content := strings.TrimSpace(payload.Content)
	if content == "" {
		writeJSONError(w, http.StatusBadRequest, "invalid_request", "Message content cannot be empty")
		return
	}

	db := sqlite.GetDB()

	// Check if recipient exists
	var recipientExists bool
	err := db.QueryRow(`SELECT EXISTS(SELECT 1 FROM Users WHERE user_id = ?)`, payload.RecipientID).Scan(&recipientExists)
	if err != nil || !recipientExists {
		writeJSONError(w, http.StatusNotFound, "not_found", "Recipient not found")
		return
	}

	// Check if sender can send message to recipient
	canSend, err := canSendMessage(db, senderID, payload.RecipientID)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "db_error", "Failed to check message permissions")
		return
	}
	if !canSend {
		writeJSONError(w, http.StatusForbidden, "forbidden", "You can only message users who follow you or have a public profile")
		return
	}

	// Insert message into database
	result, err := db.Exec(`
		INSERT INTO Messages (sender_id, recipient_id, content)
		VALUES (?, ?, ?)
	`, senderID, payload.RecipientID, content)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "db_error", "Failed to send message")
		return
	}

	messageID, _ := result.LastInsertId()

	// Get sender name for WebSocket message
	var senderName string
	db.QueryRow(`
		SELECT COALESCE(NULLIF(nickname, ''), email) 
		FROM Users WHERE user_id = ?
	`, senderID).Scan(&senderName)

	// Check if recipient can receive message via WebSocket
	canReceive, _ := canReceiveMessage(db, senderID, payload.RecipientID)
	if canReceive {
		// Send real-time message via WebSocket
		wsMessage := &websocket.Message{
			Type:      "private_message",
			UserID:    senderID,
			Recipient:  payload.RecipientID,
			Data: map[string]interface{}{
				"message_id":  int(messageID),
				"sender_id":   senderID,
				"sender_name": senderName,
				"content":     content,
				"created_at":  time.Now().Format(time.RFC3339),
			},
		}
		hub := websocket.GetGlobalHub()
		if hub != nil {
			hub.BroadcastToUser(payload.RecipientID, wsMessage)
		}
	}

	// Create notification for recipient (if not connected via WebSocket)
	_, _ = db.Exec(`
		INSERT INTO Notifications (user_id, type, related_user_id, message)
		VALUES (?, 'message', ?, ?)
	`, payload.RecipientID, senderID, "New message from "+senderName)

	// Send notification via WebSocket
	websocket.SendNotificationWithUser(payload.RecipientID, "message", "New message from "+senderName, senderID)

	writeJSON(w, http.StatusCreated, map[string]interface{}{
		"message_id": messageID,
		"status":     "sent",
	})
}

// GetMessagesHandler handles GET /api/messages
// Returns conversation history with a specific user
func GetMessagesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSONError(w, http.StatusMethodNotAllowed, "method_not_allowed", "Method not allowed")
		return
	}

	// Get authenticated user
	isAuth, userID, _ := utils.CheckAuth(r)
	if !isAuth || userID == 0 {
		writeJSONError(w, http.StatusUnauthorized, "unauthorized", "Please log in to view messages")
		return
	}

	// Get query parameters
	otherUserIDStr := r.URL.Query().Get("user_id")
	if otherUserIDStr == "" {
		writeJSONError(w, http.StatusBadRequest, "invalid_request", "user_id parameter is required")
		return
	}

	otherUserID, err := strconv.Atoi(otherUserIDStr)
	if err != nil || otherUserID <= 0 {
		writeJSONError(w, http.StatusBadRequest, "invalid_request", "Invalid user_id")
		return
	}

	limitStr := r.URL.Query().Get("limit")
	limit := 50 // default limit
	if limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 && parsedLimit <= 100 {
			limit = parsedLimit
		}
	}

	db := sqlite.GetDB()

	// Get messages between the two users
	rows, err := db.Query(`
		SELECT 
			m.id, m.sender_id, m.recipient_id, m.content, m.created_at, m.read_at,
			COALESCE(NULLIF(u.nickname, ''), u.email) AS sender_name
		FROM Messages m
		JOIN Users u ON u.user_id = m.sender_id
		WHERE 
			((m.sender_id = ? AND m.recipient_id = ?) OR 
			 (m.sender_id = ? AND m.recipient_id = ?))
			AND m.group_id IS NULL
		ORDER BY m.created_at ASC
		LIMIT ?
	`, userID, otherUserID, otherUserID, userID, limit)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "db_error", "Failed to fetch messages")
		return
	}
	defer rows.Close()

	messages := make([]MessageResponse, 0)
	for rows.Next() {
		var m MessageResponse
		var readAt sql.NullString

		err := rows.Scan(
			&m.ID,
			&m.SenderID,
			&m.RecipientID,
			&m.Content,
			&m.CreatedAt,
			&readAt,
			&m.SenderName,
		)
		if err != nil {
			continue
		}

		if readAt.Valid {
			readAtStr := readAt.String
			m.ReadAt = &readAtStr
		}

		messages = append(messages, m)
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"messages":    messages,
		"other_user_id": otherUserID,
		"total":       len(messages),
	})
}

// GetConversationsHandler handles GET /api/messages/conversations
// Returns list of users the current user has conversations with
func GetConversationsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSONError(w, http.StatusMethodNotAllowed, "method_not_allowed", "Method not allowed")
		return
	}

	isAuth, userID, _ := utils.CheckAuth(r)
	if !isAuth || userID == 0 {
		writeJSONError(w, http.StatusUnauthorized, "unauthorized", "Please log in")
		return
	}

	db := sqlite.GetDB()

	rows, err := db.Query(`
		SELECT DISTINCT
			CASE 
				WHEN m.sender_id = ? THEN m.recipient_id
				ELSE m.sender_id
			END AS other_user_id,
			COALESCE(NULLIF(u.nickname, ''), u.email) AS other_user_name,
			u.avatar_path AS other_user_avatar,
			MAX(m.created_at) AS last_message_time,
			COUNT(CASE WHEN m.read_at IS NULL AND m.recipient_id = ? THEN 1 END) AS unread_count
		FROM Messages m
		JOIN Users u ON u.user_id = CASE 
			WHEN m.sender_id = ? THEN m.recipient_id
			ELSE m.sender_id
		END
		WHERE (m.sender_id = ? OR m.recipient_id = ?) AND m.group_id IS NULL
		GROUP BY other_user_id, other_user_name, other_user_avatar
		ORDER BY last_message_time DESC
	`, userID, userID, userID, userID, userID)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "db_error", "Failed to fetch conversations")
		return
	}
	defer rows.Close()

	type Conversation struct {
		UserID         int    `json:"user_id"`
		UserName       string `json:"user_name"`
		UserAvatar     string `json:"user_avatar,omitempty"`
		LastMessageAt  string `json:"last_message_at"`
		UnreadCount    int    `json:"unread_count"`
	}

	conversations := make([]Conversation, 0)
	for rows.Next() {
		var c Conversation
		err := rows.Scan(
			&c.UserID,
			&c.UserName,
			&c.UserAvatar,
			&c.LastMessageAt,
			&c.UnreadCount,
		)
		if err != nil {
			continue
		}
		conversations = append(conversations, c)
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"conversations": conversations,
		"total":        len(conversations),
	})
}

// GetContactsHandler handles GET /api/messages/contacts
// Returns users the current user can message (accepted follow in either direction or public profile)
func GetContactsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSONError(w, http.StatusMethodNotAllowed, "method_not_allowed", "Method not allowed")
		return
	}

	isAuth, userID, _ := utils.CheckAuth(r)
	if !isAuth || userID == 0 {
		writeJSONError(w, http.StatusUnauthorized, "unauthorized", "Please log in")
		return
	}

	db := sqlite.GetDB()

	rows, err := db.Query(`
		SELECT DISTINCT u.user_id, COALESCE(NULLIF(u.nickname, ''), u.email) AS user_name, u.avatar_path AS user_avatar
		FROM (
			SELECT following_id AS uid FROM Followers WHERE follower_id = ? AND status = 'accepted'
			UNION
			SELECT follower_id AS uid FROM Followers WHERE following_id = ? AND status = 'accepted'
			UNION
			SELECT user_id AS uid FROM Users WHERE is_public = 1 AND user_id != ?
		) AS allowed
		JOIN Users u ON u.user_id = allowed.uid
		WHERE u.user_id != ?
		ORDER BY u.nickname, u.user_id
	`, userID, userID, userID, userID)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "db_error", "Failed to fetch contacts")
		return
	}
	defer rows.Close()

	type Contact struct {
		UserID     int    `json:"user_id"`
		UserName   string `json:"user_name"`
		UserAvatar string `json:"user_avatar,omitempty"`
	}

	contacts := make([]Contact, 0)
	for rows.Next() {
		var c Contact
		if err := rows.Scan(&c.UserID, &c.UserName, &c.UserAvatar); err != nil {
			continue
		}
		contacts = append(contacts, c)
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"contacts": contacts,
		"total":    len(contacts),
	})
}

// MarkReadHandler handles POST /api/messages/read
// Marks messages as read
func MarkReadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSONError(w, http.StatusMethodNotAllowed, "method_not_allowed", "Method not allowed")
		return
	}

	// Get authenticated user
	isAuth, userID, _ := utils.CheckAuth(r)
	if !isAuth || userID == 0 {
		writeJSONError(w, http.StatusUnauthorized, "unauthorized", "Please log in")
		return
	}

	var payload struct {
		MessageID  *int `json:"message_id"`  // Mark specific message as read
		SenderID   *int `json:"sender_id"`   // Mark all messages from sender as read
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid_json", "Invalid JSON body")
		return
	}

	db := sqlite.GetDB()

	if payload.MessageID != nil {
		// Mark specific message as read
		_, err := db.Exec(`
			UPDATE Messages 
			SET read_at = CURRENT_TIMESTAMP 
			WHERE id = ? AND recipient_id = ? AND read_at IS NULL
		`, *payload.MessageID, userID)
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, "db_error", "Failed to mark message as read")
			return
		}
	} else if payload.SenderID != nil {
		// Mark all messages from sender as read
		_, err := db.Exec(`
			UPDATE Messages 
			SET read_at = CURRENT_TIMESTAMP 
			WHERE sender_id = ? AND recipient_id = ? AND read_at IS NULL
		`, *payload.SenderID, userID)
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, "db_error", "Failed to mark messages as read")
			return
		}
	} else {
		writeJSONError(w, http.StatusBadRequest, "invalid_request", "Either message_id or sender_id must be provided")
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"message": "Messages marked as read",
	})
}

func writeJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

func writeJSONError(w http.ResponseWriter, status int, code string, message string) {
	writeJSON(w, status, map[string]string{
		"error":   code,
		"message": message,
	})
}

