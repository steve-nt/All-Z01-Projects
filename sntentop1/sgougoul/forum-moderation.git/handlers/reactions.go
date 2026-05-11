package handlers

import (
	"database/sql"
	"net/http"
	"strconv"

	"forum/db"
	"forum/sessions"
)

// ReactPost handles like/dislike actions for POSTS.
// Only accepts POST requests from logged-in users.
func ReactPost(w http.ResponseWriter, r *http.Request) {
	// Reject any non-POST request (audit requirement: proper HTTP handling)
	if r.Method != http.MethodPost {
		RenderError(w, r, http.StatusMethodNotAllowed, "Method not allowed.")
		return
	}

	// Ensure form values are parsed before reading them
	_ = r.ParseForm()

	// Check user session (only logged-in users can react)
	userID, ok := sessions.GetUserID(r)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Read post ID from form and validate it
	postID, err := strconv.Atoi(r.FormValue("post_id"))
	if err != nil || postID <= 0 {
		RenderError(w, r, http.StatusBadRequest, "Invalid post.")
		return
	}

	// AUDIT: reactions are allowed only when the post is visible to the current user.
	// This blocks normal users from reacting to hidden moderation-only content.
	isPrivileged := false
	if ok, err := db.IsModeratorOrAdmin(userID); err == nil {
		isPrivileged = ok
	}

	if _, err := db.GetVisiblePostByID(postID, userID, isPrivileged); err != nil {
		if err == sql.ErrNoRows {
			RenderError(w, r, http.StatusNotFound, "Post not found.")
			return
		}
		RenderError(w, r, http.StatusInternalServerError, "Could not validate post visibility.")
		return
	}

	// Reaction value must be +1 (like) or -1 (dislike)
	value, err := strconv.Atoi(r.FormValue("value"))
	if err != nil || (value != 1 && value != -1) {
		RenderError(w, r, http.StatusBadRequest, "Invalid reaction.")
		return
	}

	// Insert or update reaction in database
	if err := db.UpsertPostReaction(userID, postID, value); err != nil {
		RenderError(w, r, http.StatusInternalServerError, "Could not save reaction.")
		return
	}

	// Notify the post owner about likes/dislikes from other users.
	if ownerID, err := db.GetPostOwnerID(postID); err == nil {
		notifType := "post_liked"
		if value == -1 {
			notifType = "post_disliked"
		}
		_ = db.CreateNotification(ownerID, userID, postID, notifType)
	}

	// Redirect back to the same post page
	http.Redirect(w, r, "/post?id="+strconv.Itoa(postID), http.StatusSeeOther)
}

// ReactComment handles like/dislike actions for COMMENTS.
// Same logic as ReactPost but targets a comment instead.
func ReactComment(w http.ResponseWriter, r *http.Request) {
	// Only POST allowed
	if r.Method != http.MethodPost {
		RenderError(w, r, http.StatusMethodNotAllowed, "Method not allowed.")
		return
	}

	_ = r.ParseForm()

	// Must be logged in
	userID, ok := sessions.GetUserID(r)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Validate post ID (needed for redirect later)
	postID, err := strconv.Atoi(r.FormValue("post_id"))
	if err != nil || postID <= 0 {
		RenderError(w, r, http.StatusBadRequest, "Invalid post.")
		return
	}

	// AUDIT: comment reactions inherit post visibility rules.
	// If the parent post is hidden from the current user, reacting is denied too.
	isPrivileged := false
	if ok, err := db.IsModeratorOrAdmin(userID); err == nil {
		isPrivileged = ok
	}

	if _, err := db.GetVisiblePostByID(postID, userID, isPrivileged); err != nil {
		if err == sql.ErrNoRows {
			RenderError(w, r, http.StatusNotFound, "Post not found.")
			return
		}
		RenderError(w, r, http.StatusInternalServerError, "Could not validate post visibility.")
		return
	}

	// Validate comment ID
	commentID, err := strconv.Atoi(r.FormValue("comment_id"))
	if err != nil || commentID <= 0 {
		RenderError(w, r, http.StatusBadRequest, "Invalid comment.")
		return
	}

	// Validate reaction value
	value, err := strconv.Atoi(r.FormValue("value"))
	if err != nil || (value != 1 && value != -1) {
		RenderError(w, r, http.StatusBadRequest, "Invalid reaction.")
		return
	}

	// Store reaction in DB
	if err := db.UpsertCommentReaction(userID, commentID, value); err != nil {
		RenderError(w, r, http.StatusInternalServerError, "Could not save reaction.")
		return
	}

	// Redirect back to the same post
	http.Redirect(w, r, "/post?id="+strconv.Itoa(postID), http.StatusSeeOther)
}