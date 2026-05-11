package notifications

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"social-network/backend/pkg/db/sqlite"
	"social-network/backend/utils"
)

type NotificationResponse struct {
	ID              int    `json:"id"`
	Type            string `json:"type"`
	Message         string `json:"message"`
	RelatedUserID   *int   `json:"related_user_id,omitempty"`
	RelatedGroupID  *int   `json:"related_group_id,omitempty"`
	RelatedPostID   *int   `json:"related_post_id,omitempty"`
	RelatedCommentID *int  `json:"related_comment_id,omitempty"`
	RelatedEventID  *int   `json:"related_event_id,omitempty"`
	IsRead          bool   `json:"is_read"`
	CreatedAt       string `json:"created_at"`
}

// NotificationsHandler handles GET /api/notifications
// Returns all notifications for the authenticated user
func NotificationsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSONError(w, http.StatusMethodNotAllowed, "method_not_allowed", "Method not allowed")
		return
	}

	// Get authenticated user
	isAuth, userID, _ := utils.CheckAuth(r)
	if !isAuth || userID == 0 {
		writeJSONError(w, http.StatusUnauthorized, "unauthorized", "Please log in to view notifications")
		return
	}

	db := sqlite.GetDB()

	// Get query parameters
	unreadOnly := r.URL.Query().Get("unread_only") == "true"
	limitStr := r.URL.Query().Get("limit")
	limit := 50 // default limit
	if limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 && parsedLimit <= 100 {
			limit = parsedLimit
		}
	}

	var query string
	var args []interface{}

	if unreadOnly {
		query = `
			SELECT id, type, message, related_user_id, related_group_id, 
			       related_post_id, related_comment_id, related_event_id, 
			       is_read, created_at
			FROM Notifications
			WHERE user_id = ? AND is_read = FALSE
			ORDER BY created_at DESC
			LIMIT ?
		`
		args = []interface{}{userID, limit}
	} else {
		query = `
			SELECT id, type, message, related_user_id, related_group_id, 
			       related_post_id, related_comment_id, related_event_id, 
			       is_read, created_at
			FROM Notifications
			WHERE user_id = ?
			ORDER BY created_at DESC
			LIMIT ?
		`
		args = []interface{}{userID, limit}
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "db_error", "Failed to fetch notifications")
		return
	}
	defer rows.Close()

	notifications := make([]NotificationResponse, 0)
	for rows.Next() {
		var n NotificationResponse
		var relatedUserID, relatedGroupID, relatedPostID, relatedCommentID, relatedEventID sql.NullInt64

		err := rows.Scan(
			&n.ID,
			&n.Type,
			&n.Message,
			&relatedUserID,
			&relatedGroupID,
			&relatedPostID,
			&relatedCommentID,
			&relatedEventID,
			&n.IsRead,
			&n.CreatedAt,
		)
		if err != nil {
			continue
		}

		if relatedUserID.Valid {
			val := int(relatedUserID.Int64)
			n.RelatedUserID = &val
		}
		if relatedGroupID.Valid {
			val := int(relatedGroupID.Int64)
			n.RelatedGroupID = &val
		}
		if relatedPostID.Valid {
			val := int(relatedPostID.Int64)
			n.RelatedPostID = &val
		}
		if relatedCommentID.Valid {
			val := int(relatedCommentID.Int64)
			n.RelatedCommentID = &val
		}
		if relatedEventID.Valid {
			val := int(relatedEventID.Int64)
			n.RelatedEventID = &val
		}

		notifications = append(notifications, n)
	}

	// Get unread count
	var unreadCount int
	db.QueryRow(`SELECT COUNT(*) FROM Notifications WHERE user_id = ? AND is_read = FALSE`, userID).Scan(&unreadCount)

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"notifications": notifications,
		"unread_count":  unreadCount,
		"total":         len(notifications),
	})
}

// MarkReadHandler handles POST /api/notifications/read
// Marks one or all notifications as read
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
		NotificationID *int  `json:"notification_id"` // If provided, mark only this one as read
		MarkAllRead     bool  `json:"mark_all_read"`   // If true, mark all as read
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid_json", "Invalid JSON body")
		return
	}

	db := sqlite.GetDB()

	if payload.MarkAllRead {
		// Mark all notifications as read
		_, err := db.Exec(`UPDATE Notifications SET is_read = TRUE WHERE user_id = ?`, userID)
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, "db_error", "Failed to update notifications")
			return
		}
		writeJSON(w, http.StatusOK, map[string]interface{}{
			"message": "All notifications marked as read",
		})
	} else if payload.NotificationID != nil {
		// Mark specific notification as read
		result, err := db.Exec(`UPDATE Notifications SET is_read = TRUE WHERE id = ? AND user_id = ?`, *payload.NotificationID, userID)
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, "db_error", "Failed to update notification")
			return
		}
		rowsAffected, _ := result.RowsAffected()
		if rowsAffected == 0 {
			writeJSONError(w, http.StatusNotFound, "not_found", "Notification not found")
			return
		}
		writeJSON(w, http.StatusOK, map[string]interface{}{
			"message": "Notification marked as read",
		})
	} else {
		writeJSONError(w, http.StatusBadRequest, "invalid_request", "Either notification_id or mark_all_read must be provided")
		return
	}
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

