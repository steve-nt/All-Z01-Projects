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

// isGroupMember checks if a user is a member of a group
func isGroupMember(db *sql.DB, groupID int, userID int) (bool, error) {
	var isMember bool
	err := db.QueryRow(`
		SELECT EXISTS(
			SELECT 1 FROM Group_Members
			WHERE group_id = ? AND user_id = ?
		)
	`, groupID, userID).Scan(&isMember)
	return isMember, err
}

// SendGroupMessageHandler handles POST /api/messages/group/send
// Sends a message to a group chat
func SendGroupMessageHandler(w http.ResponseWriter, r *http.Request) {
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
		GroupID int    `json:"group_id"`
		Content string `json:"content"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid_json", "Invalid JSON body")
		return
	}

	// Validate input
	if payload.GroupID <= 0 {
		writeJSONError(w, http.StatusBadRequest, "invalid_request", "group_id is required")
		return
	}

	content := strings.TrimSpace(payload.Content)
	if content == "" {
		writeJSONError(w, http.StatusBadRequest, "invalid_request", "Message content cannot be empty")
		return
	}

	db := sqlite.GetDB()

	// Check if group exists
	var groupExists bool
	err := db.QueryRow(`SELECT EXISTS(SELECT 1 FROM Groups WHERE id = ?)`, payload.GroupID).Scan(&groupExists)
	if err != nil || !groupExists {
		writeJSONError(w, http.StatusNotFound, "not_found", "Group not found")
		return
	}

	// Check if sender is a group member
	member, err := isGroupMember(db, payload.GroupID, senderID)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "db_error", "Failed to check membership")
		return
	}
	if !member {
		writeJSONError(w, http.StatusForbidden, "forbidden", "Only group members can send messages")
		return
	}

	// Insert message into database
	result, err := db.Exec(`
		INSERT INTO Messages (sender_id, group_id, content)
		VALUES (?, ?, ?)
	`, senderID, payload.GroupID, content)
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

	var groupName string
	db.QueryRow(`SELECT name FROM Groups WHERE id = ?`, payload.GroupID).Scan(&groupName)
	if groupName == "" {
		groupName = "Group"
	}

	// Get all group members (except sender) for broadcasting
	rows, err := db.Query(`
		SELECT user_id FROM Group_Members WHERE group_id = ? AND user_id != ?
	`, payload.GroupID, senderID)
	if err == nil {
		defer rows.Close()

		wsMessage := &websocket.Message{
			Type:    "group_message",
			UserID:  senderID,
			GroupID: payload.GroupID,
			Data: map[string]interface{}{
				"message_id":  int(messageID),
				"sender_id":   senderID,
				"sender_name": senderName,
				"group_id":    payload.GroupID,
				"content":     content,
				"created_at":  time.Now().Format(time.RFC3339),
			},
		}

		hub := websocket.GetGlobalHub()
		notifMsg := "New message in " + groupName + " from " + senderName
		if hub != nil {
			for rows.Next() {
				var memberID int
				if err := rows.Scan(&memberID); err != nil {
					continue
				}
				hub.BroadcastToUser(memberID, wsMessage)
				_, _ = db.Exec(`
					INSERT INTO Notifications (user_id, type, related_user_id, related_group_id, message)
					VALUES (?, 'group_message', ?, ?, ?)
				`, memberID, senderID, payload.GroupID, notifMsg)
				websocket.SendNotificationWithGroupAndUser(memberID, "group_message", notifMsg, senderID, payload.GroupID)
			}
		}
	}

	writeJSON(w, http.StatusCreated, map[string]interface{}{
		"message_id": messageID,
		"status":     "sent",
	})
}

// GetGroupMessagesHandler handles GET /api/messages/group?group_id=X
// Returns group chat history
func GetGroupMessagesHandler(w http.ResponseWriter, r *http.Request) {
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
	groupIDStr := r.URL.Query().Get("group_id")
	if groupIDStr == "" {
		writeJSONError(w, http.StatusBadRequest, "invalid_request", "group_id parameter is required")
		return
	}

	groupID, err := strconv.Atoi(groupIDStr)
	if err != nil || groupID <= 0 {
		writeJSONError(w, http.StatusBadRequest, "invalid_request", "Invalid group_id")
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

	// Check if user is a group member
	member, err := isGroupMember(db, groupID, userID)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "db_error", "Failed to check membership")
		return
	}
	if !member {
		writeJSONError(w, http.StatusForbidden, "forbidden", "Only group members can view messages")
		return
	}

	// Get group messages
	rows, err := db.Query(`
		SELECT 
			m.id, m.sender_id, m.group_id, m.content, m.created_at,
			COALESCE(NULLIF(u.nickname, ''), u.email) AS sender_name
		FROM Messages m
		JOIN Users u ON u.user_id = m.sender_id
		WHERE m.group_id = ?
		ORDER BY m.created_at ASC
		LIMIT ?
	`, groupID, limit)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "db_error", "Failed to fetch messages")
		return
	}
	defer rows.Close()

	messages := make([]MessageResponse, 0)
	for rows.Next() {
		var m MessageResponse
		var groupIDVal sql.NullInt64

		err := rows.Scan(
			&m.ID,
			&m.SenderID,
			&groupIDVal,
			&m.Content,
			&m.CreatedAt,
			&m.SenderName,
		)
		if err != nil {
			continue
		}

		if groupIDVal.Valid {
			val := int(groupIDVal.Int64)
			m.GroupID = &val
		}

		messages = append(messages, m)
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"messages": messages,
		"group_id": groupID,
		"total":    len(messages),
	})
}

