package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"realtimeforum/internals/database"
	"realtimeforum/internals/utils"
	"sort"
	"strconv"
	"strings"
	"time"
)

// SendMessageHandler handles sending private messages
func SendMessageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Check authentication
	cookie, err := r.Cookie("session")
	if err != nil || !utils.IsValidSession(cookie.Value) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	senderID := utils.GetUserIDFromSession(cookie.Value)
	if senderID == 0 {
		http.Error(w, "Invalid session", http.StatusUnauthorized)
		return
	}

	// Parse form data
	receiverIDStr := r.FormValue("receiver_id")
	content := strings.TrimSpace(r.FormValue("content"))

	if receiverIDStr == "" || content == "" {
		http.Error(w, "Receiver ID and content are required", http.StatusBadRequest)
		return
	}

	receiverID, err := strconv.Atoi(receiverIDStr)
	if err != nil {
		http.Error(w, "Invalid receiver ID", http.StatusBadRequest)
		return
	}

	if receiverID == senderID {
		http.Error(w, "Cannot send message to yourself", http.StatusBadRequest)
		return
	}

	// Verify receiver exists
	db := database.CreateTable()
	defer db.Close()

	var receiverExists int
	err = db.QueryRow("SELECT COUNT(*) FROM Users WHERE user_id = ?", receiverID).Scan(&receiverExists)
	if err != nil || receiverExists == 0 {
		http.Error(w, "Receiver not found", http.StatusNotFound)
		return
	}

	// Insert message into database
	result, err := db.Exec(
		"INSERT INTO PrivateMessages (sender_id, receiver_id, content) VALUES (?, ?, ?)",
		senderID, receiverID, content,
	)
	if err != nil {
		http.Error(w, "Failed to send message", http.StatusInternalServerError)
		return
	}

	messageID, _ := result.LastInsertId()

	// Get sender username for WebSocket notification
	var senderUsername string
	db.QueryRow("SELECT username FROM Users WHERE user_id = ?", senderID).Scan(&senderUsername)

	// Send real-time notification via WebSocket
	// wsBroadcast is in the same package (handlers), so it's accessible
	event := map[string]interface{}{
		"type":            "private_message",
		"message_id":      messageID,
		"sender_id":       senderID,
		"sender_username": senderUsername,
		"sender_name":     senderUsername, // Keep both for compatibility
		"receiver_id":     receiverID,
		"content":         content,
		"timestamp":       time.Now().Format(time.RFC3339),
		"date":            time.Now().Format("2006-01-02"),
		"time":            time.Now().Format("15:04"),
		"datetime":        time.Now().Format("2006-01-02 15:04:05"),
	}
	if jsonEvent, err := json.Marshal(event); err == nil {
		log.Printf("Broadcasting private message: sender_id=%d, receiver_id=%d\n", senderID, receiverID)
		wsBroadcast <- jsonEvent
	} else {
		log.Printf("Failed to marshal WebSocket event: %v\n", err)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":    true,
		"message_id": messageID,
	})
}

// GetMessagesHandler returns messages between two users with pagination
func GetMessagesHandler(w http.ResponseWriter, r *http.Request) {
	// Check authentication
	cookie, err := r.Cookie("session")
	if err != nil || !utils.IsValidSession(cookie.Value) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	currentUserID := utils.GetUserIDFromSession(cookie.Value)
	if currentUserID == 0 {
		http.Error(w, "Invalid session", http.StatusUnauthorized)
		return
	}

	// Get other user ID
	otherUserIDStr := r.URL.Query().Get("user_id")
	offsetStr := r.URL.Query().Get("offset")
	limitStr := r.URL.Query().Get("limit")

	if otherUserIDStr == "" {
		http.Error(w, "User ID required", http.StatusBadRequest)
		return
	}

	otherUserID, err := strconv.Atoi(otherUserIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Default pagination
	offset := 0
	limit := 10
	if offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil {
			offset = o
		}
	}
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 50 {
			limit = l
		}
	}

	db := database.CreateTable()
	defer db.Close()

	// Get messages between the two users (bidirectional)
	query := `
		SELECT message_id, sender_id, receiver_id, content, is_read, creation_date,
		       u.username as sender_username
		FROM PrivateMessages pm
		JOIN Users u ON pm.sender_id = u.user_id
		WHERE (sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)
		ORDER BY creation_date DESC
		LIMIT ? OFFSET ?
	`

	rows, err := db.Query(query, currentUserID, otherUserID, otherUserID, currentUserID, limit, offset)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var messages []map[string]interface{}
	for rows.Next() {
		var msg map[string]interface{} = make(map[string]interface{})
		var messageID, senderID, receiverID int
		var content string
		var isRead bool
		var creationDate time.Time
		var senderUsername string

		err := rows.Scan(&messageID, &senderID, &receiverID, &content, &isRead, &creationDate, &senderUsername)
		if err != nil {
			continue
		}

		msg["id"] = messageID
		msg["sender_id"] = senderID
		msg["receiver_id"] = receiverID
		msg["content"] = content
		msg["is_read"] = isRead
		msg["sender_username"] = senderUsername
		msg["timestamp"] = creationDate.Format(time.RFC3339)
		msg["date"] = creationDate.Format("2006-01-02")
		msg["time"] = creationDate.Format("15:04")
		msg["datetime"] = creationDate.Format("2006-01-02 15:04:05")
		msg["time_ago"] = utils.FormatTimeAgo(creationDate)
		msg["is_sender"] = senderID == currentUserID

		messages = append(messages, msg)
	}

	// Mark messages as read if they were sent to current user
	if len(messages) > 0 {
		db.Exec(`
			UPDATE PrivateMessages 
			SET is_read = 1 
			WHERE receiver_id = ? AND sender_id = ? AND is_read = 0
		`, currentUserID, otherUserID)
	}

	// Reverse to show oldest first (for chat display)
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}

// GetUsersListHandler returns list of users with their last message and online status
func GetUsersListHandler(w http.ResponseWriter, r *http.Request) {
	// Check authentication
	cookie, err := r.Cookie("session")
	if err != nil || !utils.IsValidSession(cookie.Value) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	currentUserID := utils.GetUserIDFromSession(cookie.Value)
	if currentUserID == 0 {
		http.Error(w, "Invalid session", http.StatusUnauthorized)
		return
	}

	db := database.CreateTable()
	defer db.Close()

	// Get all users except current user, with their last message info
	// Note: is_online will be determined by checking WebSocket connections, not Sessions
	// Simplified query to ensure it works correctly
	query := `
		SELECT 
			u.user_id,
			u.username,
			u.first_name,
			u.last_name
		FROM Users u
		WHERE u.user_id != ?
		ORDER BY u.username ASC
	`

	rows, err := db.Query(query, currentUserID)
	if err != nil {
		log.Printf("Database error in GetUsersListHandler: %v\n", err)
		// Return empty array instead of error
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]map[string]interface{}{})
		return
	}
	defer rows.Close()

	var users []map[string]interface{}
	for rows.Next() {
		var user map[string]interface{} = make(map[string]interface{})
		var userID int
		var username, firstName, lastName string

		err := rows.Scan(&userID, &username, &firstName, &lastName)
		if err != nil {
			log.Printf("Error scanning user row: %v\n", err)
			continue
		}

		// Check if user is online via WebSocket connections
		isOnline := IsUserOnline(userID)

		// Get unread message count for this user
		var unreadCount int
		err = db.QueryRow(`
			SELECT COUNT(*) 
			FROM PrivateMessages
			WHERE sender_id = ? AND receiver_id = ? AND (is_read = 0 OR is_read IS NULL)
		`, userID, currentUserID).Scan(&unreadCount)
		if err != nil {
			unreadCount = 0 // Default to 0 if query fails
		}

		// Get last message info for this user
		var lastMessageDate sql.NullTime
		var lastMessageContent sql.NullString
		var lastMessageSenderID sql.NullInt64
		var lastMessageIsRead sql.NullBool

		err = db.QueryRow(`
			SELECT creation_date, content, sender_id, is_read
			FROM PrivateMessages
			WHERE (sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)
			ORDER BY creation_date DESC
			LIMIT 1
		`, userID, currentUserID, currentUserID, userID).Scan(
			&lastMessageDate, &lastMessageContent, &lastMessageSenderID, &lastMessageIsRead,
		)

		// Filter offline users: only include if they were active in the last hour
		if !isOnline {
			// Check if user has any activity in the last hour
			// Activity can be: last message OR recent session (not expired or expired within last hour)
			hasRecentActivity := false
			
			// Check if they have a message in the last hour
			if lastMessageDate.Valid {
				oneHourAgo := time.Now().Add(-1 * time.Hour)
				if lastMessageDate.Time.After(oneHourAgo) {
					hasRecentActivity = true
				}
			}
			
			// If no recent message, check for recent session activity
			if !hasRecentActivity {
				var lastSessionExpiration sql.NullTime
				// Get most recent session expiration date
				db.QueryRow(`
					SELECT MAX(expiration_date) 
					FROM Sessions 
					WHERE user_id = ?
				`, userID).Scan(&lastSessionExpiration)
				
				if lastSessionExpiration.Valid {
					oneHourAgo := time.Now().Add(-1 * time.Hour)
					// If session hasn't expired yet, or expired within the last hour, they were active recently
					if lastSessionExpiration.Time.After(oneHourAgo) {
						hasRecentActivity = true
					}
				}
			}
			
			// Skip offline users who weren't active in the last hour
			if !hasRecentActivity {
				continue
			}
		}

		user["id"] = userID
		user["username"] = username
		user["first_name"] = firstName
		user["last_name"] = lastName
		user["is_online"] = isOnline
		user["unread_count"] = unreadCount

		if lastMessageDate.Valid {
			user["last_message_date"] = lastMessageDate.Time.Format(time.RFC3339)
			user["last_message_time_ago"] = utils.FormatTimeAgo(lastMessageDate.Time)
		} else {
			user["last_message_date"] = nil
			user["last_message_time_ago"] = nil
		}

		if lastMessageContent.Valid {
			user["last_message_content"] = lastMessageContent.String
		} else {
			user["last_message_content"] = nil
		}

		if lastMessageSenderID.Valid {
			user["last_message_sender_id"] = lastMessageSenderID.Int64
			user["is_last_message_from_me"] = int(lastMessageSenderID.Int64) == currentUserID
		} else {
			user["last_message_sender_id"] = nil
			user["is_last_message_from_me"] = nil
		}

		if lastMessageIsRead.Valid {
			user["last_message_is_read"] = lastMessageIsRead.Bool
		} else {
			user["last_message_is_read"] = nil
		}

		users = append(users, user)
	}
	
	// Sort users by last message date (most recent first), then alphabetically
	sort.Slice(users, func(i, j int) bool {
		dateI := users[i]["last_message_date"]
		dateJ := users[j]["last_message_date"]
		
		if dateI == nil && dateJ == nil {
			// Both have no messages, sort alphabetically
			return users[i]["username"].(string) < users[j]["username"].(string)
		}
		if dateI == nil {
			return false // i has no message, j comes first
		}
		if dateJ == nil {
			return true // j has no message, i comes first
		}
		
		// Both have messages, compare dates
		dateIStr := dateI.(string)
		dateJStr := dateJ.(string)
		if dateIStr != dateJStr {
			return dateIStr > dateJStr // Most recent first
		}
		
		// Same date, sort alphabetically
		return users[i]["username"].(string) < users[j]["username"].(string)
	})

	// Ensure we always return an array, even if empty
	if users == nil {
		users = make([]map[string]interface{}, 0)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

