package profile

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"social-network/backend/pkg/db/sqlite"
	"social-network/backend/pkg/websocket"
	"social-network/backend/utils"
)

// ProfileViewHandler returns profile data with privacy-aware filtering
func ProfileViewHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userIDParam := r.URL.Query().Get("user_id")
	if userIDParam == "" {
		http.Error(w, "user_id is required", http.StatusBadRequest)
		return
	}

	ownerID, err := strconv.Atoi(userIDParam)
	if err != nil || ownerID <= 0 {
		http.Error(w, "invalid user_id", http.StatusBadRequest)
		return
	}

	isAuth, viewerID, _ := utils.CheckAuth(r)
	if !isAuth {
		viewerID = 0
	}

	db := sqlite.GetDB()

	var (
		userID             int
		email              string
		firstName          string
		lastName           string
		nickname           string
		avatar             string
		aboutMe            string
		isPublic           bool
		dateOfBirth        string
		relationshipStatus sql.NullString
		hobbies            sql.NullString
	)

	err = db.QueryRow(`
		SELECT user_id, email, first_name, last_name, nickname, avatar_path, about_me, is_public,
		       date_of_birth, relationship_status, hobbies
		FROM Users
		WHERE user_id = ?
	`, ownerID).Scan(&userID, &email, &firstName, &lastName, &nickname, &avatar, &aboutMe, &isPublic,
		&dateOfBirth, &relationshipStatus, &hobbies)
	if err != nil {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	// Allow full view for owner or public profiles
	if isPublic || (viewerID != 0 && viewerID == userID) {
		json.NewEncoder(w).Encode(buildProfileResponse(userID, email, firstName, lastName, nickname, avatar, aboutMe, isPublic, dateOfBirth, relationshipStatus, hobbies))
		return
	}
	// TODO(Part2): viewer identity comes from utils.CheckAuth(r) -> viewerID (0 if not logged in)
	// TODO(Part2): privacy gate uses Users.is_public for the target user (ownerID/userID)
	// TODO(Part2): for private profiles, allow ONLY if Followers has follower_id=viewerID AND following_id=userID AND status='accepted'
	// Private profile: require accepted follower relationship
	if viewerID == 0 {
		// Not logged in: return limited profile (name + avatar only)
		json.NewEncoder(w).Encode(buildLimitedProfileResponse(userID, firstName, lastName, nickname, avatar))
		return
	}

	var isFollower bool
	err = db.QueryRow(`
		SELECT EXISTS (
			SELECT 1 FROM Followers
			WHERE follower_id = ? AND following_id = ? AND status = 'accepted'
		)
	`, viewerID, userID).Scan(&isFollower)
	if err != nil {
		http.Error(w, "failed to check follow status", http.StatusInternalServerError)
		return
	}

	if !isFollower {
		// Return limited profile (name + avatar) so visitors see whose profile it is
		json.NewEncoder(w).Encode(buildLimitedProfileResponse(userID, firstName, lastName, nickname, avatar))
		return
	}

	json.NewEncoder(w).Encode(buildProfileResponse(userID, email, firstName, lastName, nickname, avatar, aboutMe, isPublic, dateOfBirth, relationshipStatus, hobbies))
}

func buildProfileResponse(userID int, email, firstName, lastName, nickname, avatar, aboutMe string, isPublic bool, dateOfBirth string, relationshipStatus, hobbies sql.NullString) map[string]interface{} {
	res := map[string]interface{}{
		"user_id":     userID,
		"email":       email,
		"first_name":  firstName,
		"last_name":   lastName,
		"nickname":    nickname,
		"avatar":      avatar,
		"about_me":    aboutMe,
		"is_public":   isPublic,
		"date_of_birth": dateOfBirth,
	}
	if relationshipStatus.Valid {
		res["relationship_status"] = relationshipStatus.String
	} else {
		res["relationship_status"] = ""
	}
	if hobbies.Valid {
		res["hobbies"] = hobbies.String
	} else {
		res["hobbies"] = ""
	}
	return res
}

// buildLimitedProfileResponse returns minimal profile data for private profiles when viewer has no access.
// Only safe, public-facing fields so visitors see whose profile they are viewing.
func buildLimitedProfileResponse(userID int, firstName, lastName, nickname, avatar string) map[string]interface{} {
	return map[string]interface{}{
		"user_id":     userID,
		"first_name":  firstName,
		"last_name":   lastName,
		"nickname":    nickname,
		"avatar":      avatar,
		"is_public":   false,
		"limited":     true,
	}
}

// ProfileUpdateHandler updates profile fields for the authenticated user (own profile only).
// Accepts nickname, about_me, relationship_status, hobbies, date_of_birth (all optional).
func ProfileUpdateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost && r.Method != http.MethodPut && r.Method != "PATCH" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	cookie, err := r.Cookie("session")
	if err != nil {
		writeJSONError(w, http.StatusUnauthorized, "missing session", "unauthorized")
		return
	}
	userID := utils.GetUserIDFromSession(cookie.Value)
	if userID == 0 {
		writeJSONError(w, http.StatusUnauthorized, "invalid session", "unauthorized")
		return
	}

	var payload struct {
		Nickname           *string `json:"nickname"`
		AboutMe            *string `json:"about_me"`
		RelationshipStatus *string `json:"relationship_status"`
		Hobbies            *string `json:"hobbies"`
		DateOfBirth        *string `json:"date_of_birth"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid JSON body", "invalid_json")
		return
	}

	db := sqlite.GetDB()

	// Build dynamic update: only update fields that are present in payload
	updates := []string{}
	args := []interface{}{}
	if payload.Nickname != nil {
		updates = append(updates, "nickname = ?")
		args = append(args, *payload.Nickname)
	}
	if payload.AboutMe != nil {
		updates = append(updates, "about_me = ?")
		args = append(args, *payload.AboutMe)
	}
	if payload.RelationshipStatus != nil {
		updates = append(updates, "relationship_status = ?")
		args = append(args, *payload.RelationshipStatus)
	}
	if payload.Hobbies != nil {
		updates = append(updates, "hobbies = ?")
		args = append(args, *payload.Hobbies)
	}
	if payload.DateOfBirth != nil && strings.TrimSpace(*payload.DateOfBirth) != "" {
		updates = append(updates, "date_of_birth = ?")
		args = append(args, *payload.DateOfBirth)
	}

	if len(updates) == 0 {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"updated": true, "user_id": userID})
		return
	}

	args = append(args, userID)
	query := "UPDATE Users SET " + strings.Join(updates, ", ") + " WHERE user_id = ?"
	_, err = db.Exec(query, args...)
	if err != nil {
		log.Printf("ProfileUpdateHandler: %v", err)
		writeJSONError(w, http.StatusInternalServerError, "failed to update profile", "internal_error")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"updated": true, "user_id": userID})
}

// ProfilePrivacyHandler updates the privacy setting for the authenticated user
func ProfilePrivacyHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	cookie, err := r.Cookie("session")
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "missing session", "internal_error")
		return
	}
	userID := utils.GetUserIDFromSession(cookie.Value)

	var payload struct {
		IsPublic *bool `json:"is_public"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return
	}

	if payload.IsPublic == nil {
		http.Error(w, "is_public is required", http.StatusBadRequest)
		return
	}

	db := sqlite.GetDB()

	_, err = db.Exec(`UPDATE Users SET is_public = ? WHERE user_id = ?`, *payload.IsPublic, userID)
	if err != nil {
		log.Printf("ProfilePrivacyHandler UPDATE error: %v", err)
		http.Error(w, "failed to update privacy setting", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"user_id":   userID,
		"is_public": *payload.IsPublic,
	})
}

func writeJSONError(w http.ResponseWriter, status int, message string, code string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{
		"error":   code,
		"message": message,
	})
}

// FollowRequestHandler creates a follow relationship or pending request
func FollowRequestHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSONError(w, http.StatusMethodNotAllowed, "Method not allowed", "method_not_allowed")
		return
	}

	cookie, err := r.Cookie("session")
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "missing session", "internal_error")
		return
	}
	currentUserID := utils.GetUserIDFromSession(cookie.Value)

	var payload struct {
		UserID int `json:"user_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid JSON body", "invalid_json")
		return
	}
	if payload.UserID <= 0 {
		writeJSONError(w, http.StatusBadRequest, "user_id is required", "invalid_request")
		return
	}
	if payload.UserID == currentUserID {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error":   "invalid_request",
			"message": "cannot follow yourself",
		})
		return
	}

	db := sqlite.GetDB()

	var isPublic bool
	err = db.QueryRow(`SELECT is_public FROM Users WHERE user_id = ?`, payload.UserID).Scan(&isPublic)
	if err != nil {
		if err == sql.ErrNoRows {
			writeJSONError(w, http.StatusNotFound, "target user not found", "not_found")
			return
		}
		writeJSONError(w, http.StatusInternalServerError, "failed to fetch target user", "internal_error")
		return
	}

	var exists bool
	err = db.QueryRow(`
		SELECT EXISTS (
			SELECT 1 FROM Followers WHERE follower_id = ? AND following_id = ?
		)
	`, currentUserID, payload.UserID).Scan(&exists)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "failed to check existing follow", "internal_error")
		return
	}
	if exists {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error":   "duplicate_request",
			"message": "follow relationship already exists",
		})
		return
	}

	status := "pending"
	if isPublic {
		status = "accepted"
	}

	_, err = db.Exec(`
		INSERT INTO Followers (follower_id, following_id, status)
		VALUES (?, ?, ?)
	`, currentUserID, payload.UserID, status)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "failed to create follow request", "internal_error")
		return
	}

	if status == "pending" {
		_, err = db.Exec(`
			INSERT INTO Notifications (user_id, type, related_user_id, message)
			VALUES (?, 'follow_request', ?, ?)
		`, payload.UserID, currentUserID, "New follow request")
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, "failed to create notification", "internal_error")
			return
		}
		websocket.SendNotificationWithUser(payload.UserID, "follow_request", "New follow request", currentUserID)
	} else {
		// Public profile: notify the user that someone started following them
		var followerNickname string
		_ = db.QueryRow(`SELECT COALESCE(NULLIF(nickname, ''), email) FROM Users WHERE user_id = ?`, currentUserID).Scan(&followerNickname)
		msg := "Someone started following you"
		if followerNickname != "" {
			msg = followerNickname + " started following you"
		}
		_, err = db.Exec(`
			INSERT INTO Notifications (user_id, type, related_user_id, message)
			VALUES (?, 'new_follower', ?, ?)
		`, payload.UserID, currentUserID, msg)
		if err != nil {
			log.Printf("Warning: failed to create new_follower notification: %v", err)
		} else {
			websocket.SendNotificationWithUser(payload.UserID, "new_follower", msg, currentUserID)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":   status,
		"target":   payload.UserID,
		"follower": currentUserID,
	})
}

// FollowAcceptHandler accepts a pending follow request for the current user
func FollowAcceptHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSONError(w, http.StatusMethodNotAllowed, "Method not allowed", "method_not_allowed")
		return
	}

	cookie, err := r.Cookie("session")
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "missing session", "internal_error")
		return
	}
	currentUserID := utils.GetUserIDFromSession(cookie.Value)

	var payload struct {
		UserID int `json:"user_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid JSON body", "invalid_json")
		return
	}
	if payload.UserID <= 0 {
		writeJSONError(w, http.StatusBadRequest, "user_id is required", "invalid_request")
		return
	}

	db := sqlite.GetDB()

	var currentStatus string
	err = db.QueryRow(`
		SELECT status FROM Followers WHERE follower_id = ? AND following_id = ?
	`, payload.UserID, currentUserID).Scan(&currentStatus)
	if err != nil {
		if err == sql.ErrNoRows {
			writeJSONError(w, http.StatusNotFound, "follow request not found", "not_found")
			return
		}
		writeJSONError(w, http.StatusInternalServerError, "failed to fetch follow request", "internal_error")
		return
	}

	if currentStatus != "pending" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error":   "invalid_status",
			"message": "only pending requests can be accepted",
		})
		return
	}

	_, err = db.Exec(`
		UPDATE Followers SET status = 'accepted'
		WHERE follower_id = ? AND following_id = ?
	`, payload.UserID, currentUserID)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "failed to accept follow request", "internal_error")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":          "accepted",
		"follower_id":     payload.UserID,
		"following_id":    currentUserID,
		"previous_status": currentStatus,
	})
}

// FollowDeclineHandler declines a pending follow request for the current user
func FollowDeclineHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSONError(w, http.StatusMethodNotAllowed, "Method not allowed", "method_not_allowed")
		return
	}

	cookie, err := r.Cookie("session")
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "missing session", "internal_error")
		return
	}
	currentUserID := utils.GetUserIDFromSession(cookie.Value)

	var payload struct {
		UserID int `json:"user_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid JSON body", "invalid_json")
		return
	}
	if payload.UserID <= 0 {
		writeJSONError(w, http.StatusBadRequest, "user_id is required", "invalid_request")
		return
	}

	db := sqlite.GetDB()

	var currentStatus string
	err = db.QueryRow(`
		SELECT status FROM Followers WHERE follower_id = ? AND following_id = ?
	`, payload.UserID, currentUserID).Scan(&currentStatus)
	if err != nil {
		if err == sql.ErrNoRows {
			writeJSONError(w, http.StatusNotFound, "follow request not found", "not_found")
			return
		}
		writeJSONError(w, http.StatusInternalServerError, "failed to fetch follow request", "internal_error")
		return
	}

	if currentStatus != "pending" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error":   "invalid_status",
			"message": "only pending requests can be declined",
		})
		return
	}

	// Keep the record with status=declined to preserve history
	_, err = db.Exec(`
		UPDATE Followers SET status = 'declined'
		WHERE follower_id = ? AND following_id = ?
	`, payload.UserID, currentUserID)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "failed to decline follow request", "internal_error")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":       "declined",
		"follower_id":  payload.UserID,
		"following_id": currentUserID,
	})
}

// FollowersHandler returns accepted followers for a user with privacy enforcement
func FollowersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSONError(w, http.StatusMethodNotAllowed, "Method not allowed", "method_not_allowed")
		return
	}

	userIDParam := r.URL.Query().Get("user_id")
	if userIDParam == "" {
		writeJSONError(w, http.StatusBadRequest, "user_id is required", "invalid_request")
		return
	}
	targetID, err := strconv.Atoi(userIDParam)
	if err != nil || targetID <= 0 {
		writeJSONError(w, http.StatusBadRequest, "invalid user_id", "invalid_request")
		return
	}

	isAuth, viewerID, _ := utils.CheckAuth(r)
	if !isAuth {
		viewerID = 0
	}

	db := sqlite.GetDB()

	var isPublic bool
	err = db.QueryRow(`SELECT is_public FROM Users WHERE user_id = ?`, targetID).Scan(&isPublic)
	if err != nil {
		if err == sql.ErrNoRows {
			writeJSONError(w, http.StatusNotFound, "user not found", "not_found")
			return
		}
		writeJSONError(w, http.StatusInternalServerError, "failed to fetch user", "internal_error")
		return
	}

	if !isPublic && viewerID != targetID {
		var isFollower bool
		err = db.QueryRow(`
			SELECT EXISTS (
				SELECT 1 FROM Followers
				WHERE follower_id = ? AND following_id = ? AND status = 'accepted'
			)
		`, viewerID, targetID).Scan(&isFollower)
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, "failed to check permission", "internal_error")
			return
		}
		if !isFollower {
			// Return empty list instead of 403 so the app doesn't show a global "Forbidden" banner
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"user_id":   targetID,
				"followers": []interface{}{},
			})
			return
		}
	}

	rows, err := db.Query(`
		SELECT f.follower_id, u.nickname, u.avatar_path
		FROM Followers f
		JOIN Users u ON u.user_id = f.follower_id
		WHERE f.following_id = ? AND f.status = 'accepted'
	`, targetID)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "failed to fetch followers", "internal_error")
		return
	}
	defer rows.Close()

	type follower struct {
		UserID   int    `json:"user_id"`
		Nickname string `json:"nickname"`
		Avatar   string `json:"avatar"`
	}
	var followers []follower
	for rows.Next() {
		var f follower
		if err := rows.Scan(&f.UserID, &f.Nickname, &f.Avatar); err != nil {
			writeJSONError(w, http.StatusInternalServerError, "failed to parse followers", "internal_error")
			return
		}
		followers = append(followers, f)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"user_id":   targetID,
		"followers": followers,
	})
}

// FollowingHandler returns accepted following list for a user with privacy enforcement
func FollowingHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSONError(w, http.StatusMethodNotAllowed, "Method not allowed", "method_not_allowed")
		return
	}

	userIDParam := r.URL.Query().Get("user_id")
	if userIDParam == "" {
		writeJSONError(w, http.StatusBadRequest, "user_id is required", "invalid_request")
		return
	}
	targetID, err := strconv.Atoi(userIDParam)
	if err != nil || targetID <= 0 {
		writeJSONError(w, http.StatusBadRequest, "invalid user_id", "invalid_request")
		return
	}

	isAuth, viewerID, _ := utils.CheckAuth(r)
	if !isAuth {
		viewerID = 0
	}

	db := sqlite.GetDB()

	var isPublic bool
	err = db.QueryRow(`SELECT is_public FROM Users WHERE user_id = ?`, targetID).Scan(&isPublic)
	if err != nil {
		if err == sql.ErrNoRows {
			writeJSONError(w, http.StatusNotFound, "user not found", "not_found")
			return
		}
		writeJSONError(w, http.StatusInternalServerError, "failed to fetch user", "internal_error")
		return
	}

	if !isPublic && viewerID != targetID {
		var isFollower bool
		err = db.QueryRow(`
			SELECT EXISTS (
				SELECT 1 FROM Followers
				WHERE follower_id = ? AND following_id = ? AND status = 'accepted'
			)
		`, viewerID, targetID).Scan(&isFollower)
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, "failed to check permission", "internal_error")
			return
		}
		if !isFollower {
			// Return empty list instead of 403 so the app doesn't show a global "Forbidden" banner
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"user_id":   targetID,
				"following": []interface{}{},
			})
			return
		}
	}

	rows, err := db.Query(`
		SELECT f.following_id, u.nickname, u.avatar_path
		FROM Followers f
		JOIN Users u ON u.user_id = f.following_id
		WHERE f.follower_id = ? AND f.status = 'accepted'
	`, targetID)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "failed to fetch following", "internal_error")
		return
	}
	defer rows.Close()

	type followingInfo struct {
		UserID   int    `json:"user_id"`
		Nickname string `json:"nickname"`
		Avatar   string `json:"avatar"`
	}
	var following []followingInfo
	for rows.Next() {
		var f followingInfo
		if err := rows.Scan(&f.UserID, &f.Nickname, &f.Avatar); err != nil {
			writeJSONError(w, http.StatusInternalServerError, "failed to parse following", "internal_error")
			return
		}
		following = append(following, f)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"user_id":   targetID,
		"following": following,
	})
}

// UnfollowHandler removes any follow relationship/request from the authenticated user to target
func UnfollowHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSONError(w, http.StatusMethodNotAllowed, "Method not allowed", "method_not_allowed")
		return
	}

	cookie, err := r.Cookie("session")
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "missing session", "internal_error")
		return
	}
	currentUserID := utils.GetUserIDFromSession(cookie.Value)

	var payload struct {
		UserID int `json:"user_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid JSON body", "invalid_json")
		return
	}
	if payload.UserID <= 0 {
		writeJSONError(w, http.StatusBadRequest, "user_id is required", "invalid_request")
		return
	}

	db := sqlite.GetDB()

	res, err := db.Exec(`
		DELETE FROM Followers
		WHERE follower_id = ? AND following_id = ?
	`, currentUserID, payload.UserID)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "failed to unfollow", "internal_error")
		return
	}
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		writeJSONError(w, http.StatusNotFound, "follow relationship not found", "not_found")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":   "unfollowed",
		"target":   payload.UserID,
		"follower": currentUserID,
	})
}

// FollowRequestsHandler lists pending follow requests targeting the authenticated user
func FollowRequestsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSONError(w, http.StatusMethodNotAllowed, "Method not allowed", "method_not_allowed")
		return
	}

	cookie, err := r.Cookie("session")
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "missing session", "internal_error")
		return
	}
	currentUserID := utils.GetUserIDFromSession(cookie.Value)

	db := sqlite.GetDB()

	rows, err := db.Query(`
		SELECT f.follower_id, u.nickname, u.avatar_path, f.status, f.created_at
		FROM Followers f
		JOIN Users u ON u.user_id = f.follower_id
		WHERE f.following_id = ? AND f.status = 'pending'
	`, currentUserID)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "failed to fetch follow requests", "internal_error")
		return
	}
	defer rows.Close()

	type requestInfo struct {
		UserID    int    `json:"user_id"`
		Nickname  string `json:"nickname"`
		Avatar    string `json:"avatar"`
		Status    string `json:"status"`
		CreatedAt string `json:"created_at"`
	}

	var requests []requestInfo
	for rows.Next() {
		var req requestInfo
		if err := rows.Scan(&req.UserID, &req.Nickname, &req.Avatar, &req.Status, &req.CreatedAt); err != nil {
			writeJSONError(w, http.StatusInternalServerError, "failed to parse follow requests", "internal_error")
			return
		}
		requests = append(requests, req)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"user_id":  currentUserID,
		"requests": requests,
	})
}

// UsersSearchHandler searches users by nickname/email and returns safe fields only.
// GET /api/users/search?q=<string>&limit=<int>
func UsersSearchHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSONError(w, http.StatusMethodNotAllowed, "Method not allowed", "method_not_allowed")
		return
	}

	_, currentUserID, _ := utils.CheckAuth(r)
	if currentUserID <= 0 {
		writeJSONError(w, http.StatusUnauthorized, "Unauthorized", "unauthorized")
		return
	}

	q := strings.TrimSpace(r.URL.Query().Get("q"))

	limit := 20
	limitParam := strings.TrimSpace(r.URL.Query().Get("limit"))
	if limitParam != "" {
		parsed, err := strconv.Atoi(limitParam)
		if err != nil || parsed <= 0 {
			writeJSONError(w, http.StatusBadRequest, "invalid limit", "invalid_request")
			return
		}
		limit = parsed
	}
	if limit > 50 {
		limit = 50
	}

	db := sqlite.GetDB()
	like := "%" + q + "%"

	rows, err := db.Query(`
		SELECT user_id, COALESCE(nickname, ''), COALESCE(avatar_path, ''), COALESCE(email, '')
		FROM Users
		WHERE user_id != ?
		  AND (
			? = ''
			OR COALESCE(nickname, '') LIKE ?
			OR COALESCE(email, '') LIKE ?
			OR COALESCE(first_name, '') LIKE ?
			OR COALESCE(last_name, '') LIKE ?
		  )
		ORDER BY created_at DESC, user_id DESC
		LIMIT ?
	`, currentUserID, q, like, like, like, like, limit)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "failed to search users", "internal_error")
		return
	}
	defer rows.Close()

	type userResult struct {
		UserID   int    `json:"user_id"`
		Nickname string `json:"nickname"`
		Avatar   string `json:"avatar"`
		UserName string `json:"user_name"`
	}

	results := make([]userResult, 0)
	for rows.Next() {
		var user userResult
		var email string
		if err := rows.Scan(&user.UserID, &user.Nickname, &user.Avatar, &email); err != nil {
			writeJSONError(w, http.StatusInternalServerError, "failed to parse users", "internal_error")
			return
		}
		if user.Nickname != "" {
			user.UserName = user.Nickname
		} else {
			user.UserName = email
		}
		results = append(results, user)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"users": results,
	})
}
