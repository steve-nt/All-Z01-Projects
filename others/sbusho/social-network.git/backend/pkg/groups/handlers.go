package groups

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"social-network/backend/pkg/db/sqlite"
	"social-network/backend/pkg/websocket"
	"social-network/backend/utils"
)

type groupSummary struct {
	GroupID     int       `json:"group_id"`
	CreatorID   int       `json:"creator_id"`
	GroupName   string    `json:"group_name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

type groupPostResponse struct {
	GroupPostID int       `json:"group_post_id"`
	GroupID     int       `json:"group_id"`
	UserID      int       `json:"user_id"`
	Author      string    `json:"author"`
	Content     string    `json:"content"`
	CreatedAt   time.Time `json:"created_at"`
	ImageURL    string    `json:"image_url,omitempty"`
}

type groupCommentResponse struct {
	GroupCommentID int       `json:"group_comment_id"`
	GroupPostID    int       `json:"group_post_id"`
	UserID         int       `json:"user_id"`
	Author         string    `json:"author"`
	Content        string    `json:"content"`
	CreatedAt      time.Time `json:"created_at"`
	ImageURL       string    `json:"image_url,omitempty"`
}

type groupEventResponse struct {
	EventID      int       `json:"event_id"`
	GroupID      int       `json:"group_id"`
	CreatorID    int       `json:"creator_id"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	EventDateTime time.Time `json:"event_datetime"`
	CreatedAt    time.Time `json:"created_at"`
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func writeJSONError(w http.ResponseWriter, status int, code string, message string) {
	writeJSON(w, status, map[string]string{
		"error":   code,
		"message": message,
	})
}

func currentUserID(r *http.Request) (int, bool) {
	c, err := r.Cookie("session")
	if err != nil || c == nil {
		return 0, false
	}
	if !utils.IsValidSession(c.Value) {
		return 0, false
	}
	uid := utils.GetUserIDFromSession(c.Value)
	return uid, uid > 0
}

func isGroupMember(db *sql.DB, groupID int, userID int) (bool, error) {
	var ok bool
	err := db.QueryRow(`
		SELECT EXISTS(
			SELECT 1 FROM Group_Members
			WHERE group_id = ? AND user_id = ?
		)
	`, groupID, userID).Scan(&ok)
	return ok, err
}

func groupCreatorID(db *sql.DB, groupID int) (int, error) {
	var creatorID int
	err := db.QueryRow(`SELECT creator_id FROM Groups WHERE id = ?`, groupID).Scan(&creatorID)
	return creatorID, err
}

// GroupsHandler
// GET  /api/groups
// POST /api/groups  (auth) JSON: { "group_name": "...", "description": "..." }
func GroupsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		listGroupsHandler(w, r)
	case http.MethodPost:
		createGroupHandler(w, r)
	default:
		writeJSONError(w, http.StatusMethodNotAllowed, "method_not_allowed", "Method not allowed")
	}
}

func listGroupsHandler(w http.ResponseWriter, r *http.Request) {
	db := sqlite.GetDB()
	viewerID, _ := currentUserID(r)

	rows, err := db.Query(`
		SELECT id, creator_id, group_name, description, created_at
		FROM Groups
		ORDER BY created_at DESC
	`)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "db_error", "failed to list groups")
		return
	}
	defer rows.Close()

	groups := make([]map[string]any, 0)
	for rows.Next() {
		var g groupSummary
		if err := rows.Scan(&g.GroupID, &g.CreatorID, &g.GroupName, &g.Description, &g.CreatedAt); err != nil {
			continue
		}
		member := false
		if viewerID > 0 {
			if ok, err := isGroupMember(db, g.GroupID, viewerID); err == nil {
				member = ok
			}
		}
		groups = append(groups, map[string]any{
			"group_id":    g.GroupID,
			"creator_id":  g.CreatorID,
			"group_name":  g.GroupName,
			"description": g.Description,
			"created_at":  g.CreatedAt,
			"is_member":   member,
		})
	}
	writeJSON(w, http.StatusOK, map[string]any{"groups": groups})
}

func createGroupHandler(w http.ResponseWriter, r *http.Request) {
	uid, ok := currentUserID(r)
	if !ok {
		writeJSONError(w, http.StatusUnauthorized, "unauthorized", "Unauthorized")
		return
	}

	var payload struct {
		GroupName   string `json:"group_name"`
		Description string `json:"description"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid_json", "invalid JSON body")
		return
	}
	name := strings.TrimSpace(payload.GroupName)
	desc := strings.TrimSpace(payload.Description)
	if name == "" || desc == "" {
		writeJSONError(w, http.StatusBadRequest, "invalid_request", "group_name and description are required")
		return
	}

	db := sqlite.GetDB()
	tx, err := db.Begin()
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "db_error", "failed to start transaction")
		return
	}
	defer tx.Rollback()

	res, err := tx.Exec(`INSERT INTO Groups (creator_id, group_name, description) VALUES (?, ?, ?)`, uid, name, desc)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "db_error", "failed to create group")
		return
	}
	groupID64, _ := res.LastInsertId()
	groupID := int(groupID64)

	// Creator becomes member
	_, err = tx.Exec(`INSERT INTO Group_Members (group_id, user_id, role) VALUES (?, ?, 'creator')`, groupID, uid)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "db_error", "failed to add creator as member")
		return
	}

	if err := tx.Commit(); err != nil {
		writeJSONError(w, http.StatusInternalServerError, "db_error", "failed to commit")
		return
	}

	writeJSON(w, http.StatusCreated, map[string]any{"group_id": groupID})
}

// GroupViewHandler
// GET /api/groups/view?group_id=
func GroupViewHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSONError(w, http.StatusMethodNotAllowed, "method_not_allowed", "Method not allowed")
		return
	}
	groupIDStr := r.URL.Query().Get("group_id")
	groupID, err := strconv.Atoi(groupIDStr)
	if err != nil || groupID <= 0 {
		writeJSONError(w, http.StatusBadRequest, "invalid_request", "group_id is required")
		return
	}

	db := sqlite.GetDB()

	var g groupSummary
	err = db.QueryRow(`
		SELECT id, creator_id, group_name, description, created_at
		FROM Groups WHERE id = ?
	`, groupID).Scan(&g.GroupID, &g.CreatorID, &g.GroupName, &g.Description, &g.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			writeJSONError(w, http.StatusNotFound, "not_found", "group not found")
			return
		}
		writeJSONError(w, http.StatusInternalServerError, "db_error", "failed to fetch group")
		return
	}

	// membership info (optional if logged in)
	viewerID, _ := currentUserID(r)
	member := false
	if viewerID > 0 {
		if ok, err := isGroupMember(db, groupID, viewerID); err == nil {
			member = ok
		}
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"group":     g,
		"is_member": member,
	})
}

// InviteToGroupHandler
// POST /api/groups/invite  JSON: { "group_id": 1, "user_id": 2 }
func InviteToGroupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSONError(w, http.StatusMethodNotAllowed, "method_not_allowed", "Method not allowed")
		return
	}
	uid, ok := currentUserID(r)
	if !ok {
		writeJSONError(w, http.StatusUnauthorized, "unauthorized", "Unauthorized")
		return
	}

	var payload struct {
		GroupID int `json:"group_id"`
		UserID  int `json:"user_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid_json", "invalid JSON body")
		return
	}
	if payload.GroupID <= 0 || payload.UserID <= 0 {
		writeJSONError(w, http.StatusBadRequest, "invalid_request", "group_id and user_id are required")
		return
	}
	if payload.UserID == uid {
		writeJSONError(w, http.StatusBadRequest, "invalid_request", "cannot invite yourself")
		return
	}

	db := sqlite.GetDB()
	isMember, err := isGroupMember(db, payload.GroupID, uid)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "db_error", "failed to check membership")
		return
	}
	if !isMember {
		writeJSONError(w, http.StatusForbidden, "forbidden", "only group members can invite")
		return
	}

	// Disallow inviting existing members
	targetIsMember, err := isGroupMember(db, payload.GroupID, payload.UserID)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "db_error", "failed to check target membership")
		return
	}
	if targetIsMember {
		writeJSONError(w, http.StatusBadRequest, "invalid_request", "user is already a member")
		return
	}

	tx, err := db.Begin()
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "db_error", "failed to start transaction")
		return
	}
	defer tx.Rollback()

	res, err := tx.Exec(`
		INSERT INTO Group_Invitations (group_id, inviter_id, invitee_id, status)
		VALUES (?, ?, ?, 'pending')
	`, payload.GroupID, uid, payload.UserID)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid_request", "invitation already exists or invalid")
		return
	}
	invID, _ := res.LastInsertId()

	// Notification to invitee
	_, _ = tx.Exec(`
		INSERT INTO Notifications (user_id, type, related_user_id, related_group_id, message)
		VALUES (?, 'group_invitation', ?, ?, ?)
	`, payload.UserID, uid, payload.GroupID, "You have been invited to a group")

	if err := tx.Commit(); err != nil {
		writeJSONError(w, http.StatusInternalServerError, "db_error", "failed to commit")
		return
	}
	
	// Send real-time notification via WebSocket
	websocket.SendNotificationWithGroupAndUser(payload.UserID, "group_invitation", "You have been invited to a group", uid, payload.GroupID)

	writeJSON(w, http.StatusCreated, map[string]any{
		"invitation_id": invID,
		"status":        "pending",
	})
}

// RespondGroupInviteHandler
// POST /api/groups/invitations/respond  JSON: { "invitation_id": 1, "response": "accepted"|"declined" }
func RespondGroupInviteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSONError(w, http.StatusMethodNotAllowed, "method_not_allowed", "Method not allowed")
		return
	}
	uid, ok := currentUserID(r)
	if !ok {
		writeJSONError(w, http.StatusUnauthorized, "unauthorized", "Unauthorized")
		return
	}

	var payload struct {
		InvitationID int    `json:"invitation_id"`
		Response     string `json:"response"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid_json", "invalid JSON body")
		return
	}
	if payload.InvitationID <= 0 {
		writeJSONError(w, http.StatusBadRequest, "invalid_request", "invitation_id is required")
		return
	}
	resp := strings.ToLower(strings.TrimSpace(payload.Response))
	if resp != "accepted" && resp != "declined" {
		writeJSONError(w, http.StatusBadRequest, "invalid_request", "response must be accepted|declined")
		return
	}

	db := sqlite.GetDB()
	var groupID, inviterID int
	var status string
	err := db.QueryRow(`
		SELECT group_id, inviter_id, status
		FROM Group_Invitations
		WHERE id = ? AND invitee_id = ?
	`, payload.InvitationID, uid).Scan(&groupID, &inviterID, &status)
	if err != nil {
		if err == sql.ErrNoRows {
			writeJSONError(w, http.StatusNotFound, "not_found", "invitation not found")
			return
		}
		writeJSONError(w, http.StatusInternalServerError, "db_error", "failed to fetch invitation")
		return
	}
	if status != "pending" {
		writeJSONError(w, http.StatusBadRequest, "invalid_request", "invitation is not pending")
		return
	}

	tx, err := db.Begin()
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "db_error", "failed to start transaction")
		return
	}
	defer tx.Rollback()

	_, err = tx.Exec(`UPDATE Group_Invitations SET status = ? WHERE id = ?`, resp, payload.InvitationID)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "db_error", "failed to update invitation")
		return
	}

	if resp == "accepted" {
		// Add member (idempotent-ish; unique constraint may error if already member)
		_, _ = tx.Exec(`INSERT OR IGNORE INTO Group_Members (group_id, user_id, role) VALUES (?, ?, 'member')`, groupID, uid)
	}

	// Notify inviter (optional)
	_, _ = tx.Exec(`
		INSERT INTO Notifications (user_id, type, related_user_id, related_group_id, message)
		VALUES (?, 'group_invitation_response', ?, ?, ?)
	`, inviterID, uid, groupID, "Group invitation response: "+resp)

	if err := tx.Commit(); err != nil {
		writeJSONError(w, http.StatusInternalServerError, "db_error", "failed to commit")
		return
	}
	
	// Send real-time notification via WebSocket
	websocket.SendNotificationWithGroupAndUser(inviterID, "group_invitation_response", "Group invitation response: "+resp, uid, groupID)

	writeJSON(w, http.StatusOK, map[string]any{
		"status": resp,
	})
}

// RequestJoinGroupHandler
// POST /api/groups/join/request  JSON: { "group_id": 1 }
func RequestJoinGroupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSONError(w, http.StatusMethodNotAllowed, "method_not_allowed", "Method not allowed")
		return
	}
	uid, ok := currentUserID(r)
	if !ok {
		writeJSONError(w, http.StatusUnauthorized, "unauthorized", "Unauthorized")
		return
	}
	var payload struct {
		GroupID int `json:"group_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid_json", "invalid JSON body")
		return
	}
	if payload.GroupID <= 0 {
		writeJSONError(w, http.StatusBadRequest, "invalid_request", "group_id is required")
		return
	}

	db := sqlite.GetDB()
	if ok, err := isGroupMember(db, payload.GroupID, uid); err != nil {
		writeJSONError(w, http.StatusInternalServerError, "db_error", "failed to check membership")
		return
	} else if ok {
		writeJSONError(w, http.StatusBadRequest, "invalid_request", "already a member")
		return
	}

	creatorID, err := groupCreatorID(db, payload.GroupID)
	if err != nil {
		if err == sql.ErrNoRows {
			writeJSONError(w, http.StatusNotFound, "not_found", "group not found")
			return
		}
		writeJSONError(w, http.StatusInternalServerError, "db_error", "failed to fetch group")
		return
	}

	tx, err := db.Begin()
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "db_error", "failed to start transaction")
		return
	}
	defer tx.Rollback()

	res, err := tx.Exec(`
		INSERT INTO Group_Join_Requests (group_id, requester_id, status)
		VALUES (?, ?, 'pending')
	`, payload.GroupID, uid)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid_request", "join request already exists or invalid")
		return
	}
	reqID, _ := res.LastInsertId()

	_, _ = tx.Exec(`
		INSERT INTO Notifications (user_id, type, related_user_id, related_group_id, message)
		VALUES (?, 'group_join_request', ?, ?, ?)
	`, creatorID, uid, payload.GroupID, "New request to join your group")

	if err := tx.Commit(); err != nil {
		writeJSONError(w, http.StatusInternalServerError, "db_error", "failed to commit")
		return
	}
	
	// Send real-time notification via WebSocket
	websocket.SendNotificationWithGroupAndUser(creatorID, "group_join_request", "New request to join your group", uid, payload.GroupID)

	writeJSON(w, http.StatusCreated, map[string]any{"request_id": reqID, "status": "pending"})
}

// ListJoinRequestsHandler returns pending join requests for a group (creator only).
// GET /api/groups/join/requests?group_id=
func ListJoinRequestsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSONError(w, http.StatusMethodNotAllowed, "method_not_allowed", "Method not allowed")
		return
	}
	uid, ok := currentUserID(r)
	if !ok {
		writeJSONError(w, http.StatusUnauthorized, "unauthorized", "Unauthorized")
		return
	}
	groupID, err := parseIntQuery(r, "group_id")
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid_request", "group_id is required")
		return
	}
	db := sqlite.GetDB()
	creatorID, err := groupCreatorID(db, groupID)
	if err != nil {
		if err == sql.ErrNoRows {
			writeJSONError(w, http.StatusNotFound, "not_found", "group not found")
			return
		}
		writeJSONError(w, http.StatusInternalServerError, "db_error", "failed to fetch group")
		return
	}
	if creatorID != uid {
		writeJSONError(w, http.StatusForbidden, "forbidden", "only the group creator can view join requests")
		return
	}
	rows, err := db.Query(`
		SELECT gjr.id, gjr.requester_id,
			COALESCE(NULLIF(u.nickname, ''), u.email, '') AS requester_name
		FROM Group_Join_Requests gjr
		LEFT JOIN Users u ON u.user_id = gjr.requester_id
		WHERE gjr.group_id = ? AND gjr.status = 'pending'
		ORDER BY gjr.id ASC
	`, groupID)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "db_error", "failed to list join requests")
		return
	}
	defer rows.Close()
	type joinReqItem struct {
		RequestID     int    `json:"request_id"`
		RequesterID   int    `json:"requester_id"`
		RequesterName string `json:"requester_name"`
	}
	list := make([]joinReqItem, 0)
	for rows.Next() {
		var item joinReqItem
		var name sql.NullString
		if err := rows.Scan(&item.RequestID, &item.RequesterID, &name); err != nil {
			continue
		}
		if name.Valid {
			item.RequesterName = name.String
		}
		list = append(list, item)
	}
	writeJSON(w, http.StatusOK, map[string]any{"join_requests": list})
}

// RespondJoinRequestHandler
// POST /api/groups/join/respond  JSON: { "request_id": 1, "response": "accepted"|"declined" }
func RespondJoinRequestHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSONError(w, http.StatusMethodNotAllowed, "method_not_allowed", "Method not allowed")
		return
	}
	uid, ok := currentUserID(r)
	if !ok {
		writeJSONError(w, http.StatusUnauthorized, "unauthorized", "Unauthorized")
		return
	}

	var payload struct {
		RequestID int    `json:"request_id"`
		Response  string `json:"response"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid_json", "invalid JSON body")
		return
	}
	if payload.RequestID <= 0 {
		writeJSONError(w, http.StatusBadRequest, "invalid_request", "request_id is required")
		return
	}
	resp := strings.ToLower(strings.TrimSpace(payload.Response))
	if resp != "accepted" && resp != "declined" {
		writeJSONError(w, http.StatusBadRequest, "invalid_request", "response must be accepted|declined")
		return
	}

	db := sqlite.GetDB()
	var groupID, requesterID int
	var status string
	err := db.QueryRow(`
		SELECT group_id, requester_id, status
		FROM Group_Join_Requests
		WHERE id = ?
	`, payload.RequestID).Scan(&groupID, &requesterID, &status)
	if err != nil {
		if err == sql.ErrNoRows {
			writeJSONError(w, http.StatusNotFound, "not_found", "join request not found")
			return
		}
		writeJSONError(w, http.StatusInternalServerError, "db_error", "failed to fetch join request")
		return
	}
	if status != "pending" {
		writeJSONError(w, http.StatusBadRequest, "invalid_request", "join request is not pending")
		return
	}

	creatorID, err := groupCreatorID(db, groupID)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "db_error", "failed to fetch group")
		return
	}
	if creatorID != uid {
		writeJSONError(w, http.StatusForbidden, "forbidden", "only group creator can respond to join requests")
		return
	}

	tx, err := db.Begin()
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "db_error", "failed to start transaction")
		return
	}
	defer tx.Rollback()

	_, err = tx.Exec(`UPDATE Group_Join_Requests SET status = ? WHERE id = ?`, resp, payload.RequestID)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "db_error", "failed to update join request")
		return
	}

	if resp == "accepted" {
		_, _ = tx.Exec(`INSERT OR IGNORE INTO Group_Members (group_id, user_id, role) VALUES (?, ?, 'member')`, groupID, requesterID)
	}

	_, _ = tx.Exec(`
		INSERT INTO Notifications (user_id, type, related_user_id, related_group_id, message)
		VALUES (?, 'group_join_response', ?, ?, ?)
	`, requesterID, uid, groupID, "Group join request response: "+resp)

	if err := tx.Commit(); err != nil {
		writeJSONError(w, http.StatusInternalServerError, "db_error", "failed to commit")
		return
	}
	
	// Send real-time notification via WebSocket
	websocket.SendNotificationWithGroupAndUser(requesterID, "group_join_response", "Group join request response: "+resp, uid, groupID)

	writeJSON(w, http.StatusOK, map[string]any{"status": resp})
}

// GroupPostsHandler
// GET  /api/groups/posts?group_id=
// POST /api/groups/posts  JSON: { "group_id": 1, "content": "...", "image_url": "/frontend/uploads/images/x.gif" }
func GroupPostsHandler(w http.ResponseWriter, r *http.Request) {
	db := sqlite.GetDB()
	uid, ok := currentUserID(r)
	if !ok {
		writeJSONError(w, http.StatusUnauthorized, "unauthorized", "Unauthorized")
		return
	}

	switch r.Method {
	case http.MethodGet:
		groupID, err := parseIntQuery(r, "group_id")
		if err != nil {
			writeJSONError(w, http.StatusBadRequest, "invalid_request", "group_id is required")
			return
		}
		member, err := isGroupMember(db, groupID, uid)
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, "db_error", "failed to check membership")
			return
		}
		if !member {
			writeJSONError(w, http.StatusForbidden, "forbidden", "only members can view group posts")
			return
		}

		rows, err := db.Query(`
			SELECT
				gp.id, gp.group_id, gp.user_id,
				COALESCE(NULLIF(u.nickname, ''), u.email) AS author,
				gp.content, gp.created_at,
				COALESCE((
					SELECT gpi.image_path
					FROM Group_Posts_Images gpi
					WHERE gpi.group_post_id = gp.id
					ORDER BY gpi.id ASC
					LIMIT 1
				), '') AS image_url
			FROM Group_Posts gp
			JOIN Users u ON u.user_id = gp.user_id
			WHERE gp.group_id = ?
			ORDER BY gp.created_at DESC
		`, groupID)
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, "db_error", "failed to fetch group posts")
			return
		}
		defer rows.Close()

		posts := make([]groupPostResponse, 0)
		for rows.Next() {
			var p groupPostResponse
			var img string
			if err := rows.Scan(&p.GroupPostID, &p.GroupID, &p.UserID, &p.Author, &p.Content, &p.CreatedAt, &img); err != nil {
				continue
			}
			if img != "" {
				p.ImageURL = img
			}
			posts = append(posts, p)
		}
		writeJSON(w, http.StatusOK, map[string]any{"group_id": groupID, "posts": posts})

	case http.MethodPost:
		var payload struct {
			GroupID  int    `json:"group_id"`
			Content  string `json:"content"`
			ImageURL string `json:"image_url"`
		}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			writeJSONError(w, http.StatusBadRequest, "invalid_json", "invalid JSON body")
			return
		}
		if payload.GroupID <= 0 {
			writeJSONError(w, http.StatusBadRequest, "invalid_request", "group_id is required")
			return
		}
		content := strings.TrimSpace(payload.Content)
		hasImage := strings.TrimSpace(payload.ImageURL) != ""
		if content == "" && !hasImage {
			writeJSONError(w, http.StatusBadRequest, "invalid_request", "post must have text or an image")
			return
		}
		member, err := isGroupMember(db, payload.GroupID, uid)
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, "db_error", "failed to check membership")
			return
		}
		if !member {
			writeJSONError(w, http.StatusForbidden, "forbidden", "only members can post in the group")
			return
		}

		tx, err := db.Begin()
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, "db_error", "failed to start transaction")
			return
		}
		defer tx.Rollback()

		res, err := tx.Exec(`INSERT INTO Group_Posts (group_id, user_id, content) VALUES (?, ?, ?)`, payload.GroupID, uid, content)
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, "db_error", "failed to create group post")
			return
		}
		gpID64, _ := res.LastInsertId()
		gpID := int(gpID64)

		if strings.TrimSpace(payload.ImageURL) != "" {
			img := strings.TrimSpace(payload.ImageURL)
			imgType := inferImageTypeFromURL(img)
			if imgType == "" {
				writeJSONError(w, http.StatusBadRequest, "invalid_request", "image_url must be JPEG|PNG|GIF")
				return
			}
			_, err := tx.Exec(`INSERT INTO Group_Posts_Images (group_post_id, image_path, image_type) VALUES (?, ?, ?)`, gpID, img, imgType)
			if err != nil {
				writeJSONError(w, http.StatusInternalServerError, "db_error", "failed to save group post image")
				return
			}
		}

		if err := tx.Commit(); err != nil {
			writeJSONError(w, http.StatusInternalServerError, "db_error", "failed to commit")
			return
		}

		writeJSON(w, http.StatusCreated, map[string]any{"group_post_id": gpID})
	default:
		writeJSONError(w, http.StatusMethodNotAllowed, "method_not_allowed", "Method not allowed")
	}
}

// GroupCommentsHandler
// GET  /api/groups/comments?group_post_id=
// POST /api/groups/comments JSON: { "group_post_id": 1, "content": "...", "image_url": "/frontend/uploads/images/x.png" }
func GroupCommentsHandler(w http.ResponseWriter, r *http.Request) {
	db := sqlite.GetDB()
	uid, ok := currentUserID(r)
	if !ok {
		writeJSONError(w, http.StatusUnauthorized, "unauthorized", "Unauthorized")
		return
	}

	switch r.Method {
	case http.MethodGet:
		gpID, err := parseIntQuery(r, "group_post_id")
		if err != nil {
			writeJSONError(w, http.StatusBadRequest, "invalid_request", "group_post_id is required")
			return
		}

		// membership check via group_post -> group_id
		var groupID int
		err = db.QueryRow(`SELECT group_id FROM Group_Posts WHERE id = ?`, gpID).Scan(&groupID)
		if err != nil {
			if err == sql.ErrNoRows {
				writeJSONError(w, http.StatusNotFound, "not_found", "group post not found")
				return
			}
			writeJSONError(w, http.StatusInternalServerError, "db_error", "failed to fetch group post")
			return
		}
		member, err := isGroupMember(db, groupID, uid)
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, "db_error", "failed to check membership")
			return
		}
		if !member {
			writeJSONError(w, http.StatusForbidden, "forbidden", "only members can view comments")
			return
		}

		rows, err := db.Query(`
			SELECT
				gc.id, gc.group_post_id, gc.user_id,
				COALESCE(NULLIF(u.nickname, ''), u.email) AS author,
				gc.content, gc.created_at,
				COALESCE((
					SELECT gci.image_path
					FROM Group_Comments_Images gci
					WHERE gci.group_comment_id = gc.id
					ORDER BY gci.id ASC
					LIMIT 1
				), '') AS image_url
			FROM Group_Comments gc
			JOIN Users u ON u.user_id = gc.user_id
			WHERE gc.group_post_id = ?
			ORDER BY gc.created_at ASC
		`, gpID)
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, "db_error", "failed to fetch comments")
			return
		}
		defer rows.Close()

		out := make([]groupCommentResponse, 0)
		for rows.Next() {
			var c groupCommentResponse
			var img string
			if err := rows.Scan(&c.GroupCommentID, &c.GroupPostID, &c.UserID, &c.Author, &c.Content, &c.CreatedAt, &img); err != nil {
				continue
			}
			if img != "" {
				c.ImageURL = img
			}
			out = append(out, c)
		}
		writeJSON(w, http.StatusOK, map[string]any{"group_post_id": gpID, "comments": out})

	case http.MethodPost:
		var payload struct {
			GroupPostID int    `json:"group_post_id"`
			Content     string `json:"content"`
			ImageURL    string `json:"image_url"`
		}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			writeJSONError(w, http.StatusBadRequest, "invalid_json", "invalid JSON body")
			return
		}
		if payload.GroupPostID <= 0 {
			writeJSONError(w, http.StatusBadRequest, "invalid_request", "group_post_id is required")
			return
		}
		content := strings.TrimSpace(payload.Content)
		if content == "" {
			writeJSONError(w, http.StatusBadRequest, "invalid_request", "content is required")
			return
		}

		var groupID int
		err := db.QueryRow(`SELECT group_id FROM Group_Posts WHERE id = ?`, payload.GroupPostID).Scan(&groupID)
		if err != nil {
			if err == sql.ErrNoRows {
				writeJSONError(w, http.StatusNotFound, "not_found", "group post not found")
				return
			}
			writeJSONError(w, http.StatusInternalServerError, "db_error", "failed to fetch group post")
			return
		}

		member, err := isGroupMember(db, groupID, uid)
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, "db_error", "failed to check membership")
			return
		}
		if !member {
			writeJSONError(w, http.StatusForbidden, "forbidden", "only members can comment")
			return
		}

		tx, err := db.Begin()
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, "db_error", "failed to start transaction")
			return
		}
		defer tx.Rollback()

		res, err := tx.Exec(`INSERT INTO Group_Comments (group_post_id, user_id, content) VALUES (?, ?, ?)`, payload.GroupPostID, uid, content)
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, "db_error", "failed to create comment")
			return
		}
		cID64, _ := res.LastInsertId()
		cID := int(cID64)

		if strings.TrimSpace(payload.ImageURL) != "" {
			img := strings.TrimSpace(payload.ImageURL)
			imgType := inferImageTypeFromURL(img)
			if imgType == "" {
				writeJSONError(w, http.StatusBadRequest, "invalid_request", "image_url must be JPEG|PNG|GIF")
				return
			}
			_, err := tx.Exec(`INSERT INTO Group_Comments_Images (group_comment_id, image_path, image_type) VALUES (?, ?, ?)`, cID, img, imgType)
			if err != nil {
				writeJSONError(w, http.StatusInternalServerError, "db_error", "failed to save comment image")
				return
			}
		}

		if err := tx.Commit(); err != nil {
			writeJSONError(w, http.StatusInternalServerError, "db_error", "failed to commit")
			return
		}
		writeJSON(w, http.StatusCreated, map[string]any{"group_comment_id": cID})
	default:
		writeJSONError(w, http.StatusMethodNotAllowed, "method_not_allowed", "Method not allowed")
	}
}

// GroupEventsHandler
// GET  /api/groups/events?group_id=
// POST /api/groups/events JSON: { "group_id": 1, "title": "...", "description": "...", "event_datetime": "2026-01-22T20:30:00Z" }
func GroupEventsHandler(w http.ResponseWriter, r *http.Request) {
	db := sqlite.GetDB()
	uid, ok := currentUserID(r)
	if !ok {
		writeJSONError(w, http.StatusUnauthorized, "unauthorized", "Unauthorized")
		return
	}

	switch r.Method {
	case http.MethodGet:
		groupID, err := parseIntQuery(r, "group_id")
		if err != nil {
			writeJSONError(w, http.StatusBadRequest, "invalid_request", "group_id is required")
			return
		}
		member, err := isGroupMember(db, groupID, uid)
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, "db_error", "failed to check membership")
			return
		}
		if !member {
			writeJSONError(w, http.StatusForbidden, "forbidden", "only members can view events")
			return
		}

		rows, err := db.Query(`
			SELECT id, group_id, creator_id, title, description, event_datetime, created_at
			FROM Group_Events
			WHERE group_id = ?
			ORDER BY event_datetime ASC
		`, groupID)
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, "db_error", "failed to list events")
			return
		}
		defer rows.Close()

		events := make([]groupEventResponse, 0)
		for rows.Next() {
			var e groupEventResponse
			if err := rows.Scan(&e.EventID, &e.GroupID, &e.CreatorID, &e.Title, &e.Description, &e.EventDateTime, &e.CreatedAt); err != nil {
				continue
			}
			events = append(events, e)
		}
		writeJSON(w, http.StatusOK, map[string]any{"group_id": groupID, "events": events})

	case http.MethodPost:
		var payload struct {
			GroupID      int    `json:"group_id"`
			Title        string `json:"title"`
			Description  string `json:"description"`
			EventDateTime string `json:"event_datetime"`
		}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			writeJSONError(w, http.StatusBadRequest, "invalid_json", "invalid JSON body")
			return
		}
		if payload.GroupID <= 0 {
			writeJSONError(w, http.StatusBadRequest, "invalid_request", "group_id is required")
			return
		}
		title := strings.TrimSpace(payload.Title)
		desc := strings.TrimSpace(payload.Description)
		if title == "" || desc == "" {
			writeJSONError(w, http.StatusBadRequest, "invalid_request", "title and description are required")
			return
		}
		dt, err := time.Parse(time.RFC3339, payload.EventDateTime)
		if err != nil {
			writeJSONError(w, http.StatusBadRequest, "invalid_request", "event_datetime must be RFC3339")
			return
		}

		member, err := isGroupMember(db, payload.GroupID, uid)
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, "db_error", "failed to check membership")
			return
		}
		if !member {
			writeJSONError(w, http.StatusForbidden, "forbidden", "only members can create events")
			return
		}

		tx, err := db.Begin()
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, "db_error", "failed to start transaction")
			return
		}
		defer tx.Rollback()

		res, err := tx.Exec(`
			INSERT INTO Group_Events (group_id, creator_id, title, description, event_datetime)
			VALUES (?, ?, ?, ?, ?)
		`, payload.GroupID, uid, title, desc, dt)
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, "db_error", "failed to create event")
			return
		}
		eventID64, _ := res.LastInsertId()
		eventID := int(eventID64)

		// Notify group members (best-effort)
		rows, err := tx.Query(`SELECT user_id FROM Group_Members WHERE group_id = ?`, payload.GroupID)
		if err == nil {
			defer rows.Close()
			for rows.Next() {
				var memberID int
				if err := rows.Scan(&memberID); err != nil {
					continue
				}
				if memberID == uid {
					continue
				}
				_, _ = tx.Exec(`
					INSERT INTO Notifications (user_id, type, related_user_id, related_group_id, related_event_id, message)
					VALUES (?, 'group_event_created', ?, ?, ?, ?)
				`, memberID, uid, payload.GroupID, eventID, "New group event: "+title)
			}
		}

		if err := tx.Commit(); err != nil {
			writeJSONError(w, http.StatusInternalServerError, "db_error", "failed to commit")
			return
		}
		
		// Send real-time notifications via WebSocket to all group members
		rows2, err := db.Query(`SELECT user_id FROM Group_Members WHERE group_id = ?`, payload.GroupID)
		if err == nil {
			defer rows2.Close()
			for rows2.Next() {
				var memberID int
				if err := rows2.Scan(&memberID); err != nil {
					continue
				}
				if memberID == uid {
					continue
				}
				websocket.SendNotificationWithEvent(memberID, "group_event_created", "New group event: "+title, uid, payload.GroupID, eventID)
			}
		}
		
		writeJSON(w, http.StatusCreated, map[string]any{"event_id": eventID})
	default:
		writeJSONError(w, http.StatusMethodNotAllowed, "method_not_allowed", "Method not allowed")
	}
}

// RespondEventHandler
// POST /api/groups/events/respond JSON: { "event_id": 1, "response": "going"|"not going" }
func RespondEventHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSONError(w, http.StatusMethodNotAllowed, "method_not_allowed", "Method not allowed")
		return
	}
	uid, ok := currentUserID(r)
	if !ok {
		writeJSONError(w, http.StatusUnauthorized, "unauthorized", "Unauthorized")
		return
	}

	var payload struct {
		EventID  int    `json:"event_id"`
		Response string `json:"response"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid_json", "invalid JSON body")
		return
	}
	if payload.EventID <= 0 {
		writeJSONError(w, http.StatusBadRequest, "invalid_request", "event_id is required")
		return
	}
	resp := strings.ToLower(strings.TrimSpace(payload.Response))
	if resp != "going" && resp != "not going" {
		writeJSONError(w, http.StatusBadRequest, "invalid_request", "response must be going|not going")
		return
	}

	db := sqlite.GetDB()
	var groupID, creatorID int
	err := db.QueryRow(`SELECT group_id, creator_id FROM Group_Events WHERE id = ?`, payload.EventID).Scan(&groupID, &creatorID)
	if err != nil {
		if err == sql.ErrNoRows {
			writeJSONError(w, http.StatusNotFound, "not_found", "event not found")
			return
		}
		writeJSONError(w, http.StatusInternalServerError, "db_error", "failed to fetch event")
		return
	}

	member, err := isGroupMember(db, groupID, uid)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "db_error", "failed to check membership")
		return
	}
	if !member {
		writeJSONError(w, http.StatusForbidden, "forbidden", "only members can respond to events")
		return
	}

	tx, err := db.Begin()
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "db_error", "failed to start transaction")
		return
	}
	defer tx.Rollback()

	_, err = tx.Exec(`
		INSERT INTO Group_Event_Responses (group_event_id, user_id, response)
		VALUES (?, ?, ?)
		ON CONFLICT(group_event_id, user_id) DO UPDATE SET response = excluded.response, updated_at = CURRENT_TIMESTAMP
	`, payload.EventID, uid, resp)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "db_error", "failed to save response")
		return
	}

	if creatorID != uid {
		_, _ = tx.Exec(`
			INSERT INTO Notifications (user_id, type, related_user_id, related_group_id, related_event_id, message)
			VALUES (?, 'group_event_response', ?, ?, ?, ?)
		`, creatorID, uid, groupID, payload.EventID, "Event response: "+resp)
	}

	if err := tx.Commit(); err != nil {
		writeJSONError(w, http.StatusInternalServerError, "db_error", "failed to commit")
		return
	}
	
	// Send real-time notification via WebSocket
	if creatorID != uid {
		websocket.SendNotificationWithEvent(creatorID, "group_event_response", "Event response: "+resp, uid, groupID, payload.EventID)
	}

	writeJSON(w, http.StatusOK, map[string]any{"status": resp})
}

func parseIntQuery(r *http.Request, key string) (int, error) {
	s := r.URL.Query().Get(key)
	n, err := strconv.Atoi(s)
	if err != nil || n <= 0 {
		return 0, errors.New("invalid int query")
	}
	return n, nil
}

func filepathExt(name string) string {
	i := strings.LastIndexByte(name, '.')
	if i < 0 {
		return ""
	}
	return name[i:]
}

func inferImageTypeFromURL(u string) string {
	ext := strings.ToLower(filepathExt(u))
	switch ext {
	case ".jpg", ".jpeg":
		return "JPEG"
	case ".png":
		return "PNG"
	case ".gif":
		return "GIF"
	default:
		return ""
	}
}

