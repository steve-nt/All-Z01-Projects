package handlers

import (
	"database/sql"
	"net/http"
	"strconv"
	"strings"

	"forum/db"
	"forum/sessions"
)

// CreateComment handles posting a new comment to a post.
// Only logged-in users are allowed to comment.
func CreateComment(w http.ResponseWriter, r *http.Request) {
	// Only POST requests are allowed.
	if r.Method != http.MethodPost {
		RenderError(w, r, http.StatusMethodNotAllowed, "Method not allowed.")
		return
	}

	// Verify user session (must be logged in).
	userID, ok := sessions.GetUserID(r)
	if !ok {

		// Redirect anonymous users to login page.
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Parse post ID from submitted form.
	postID, err := strconv.Atoi(r.FormValue("post_id"))
	if err != nil || postID <= 0 {

		RenderError(w, r, http.StatusBadRequest, "Invalid post.")
		return
	}

	// AUDIT: comments may only be added to posts visible to the current user.
	// This prevents normal users from interacting with hidden moderation content.
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

	// Trim comment text to avoid empty or whitespace-only comments.
	content := strings.TrimSpace(r.FormValue("content"))
	if content == "" {

		RenderError(w, r, http.StatusBadRequest, "Comment cannot be empty.")
		return
	}

	// Insert comment into database.
	if err := db.CreateComment(postID, userID, content); err != nil {

		RenderError(w, r, http.StatusInternalServerError, "Could not save comment.")
		return
	}

	// Notify the post owner when someone else comments on their post.
	// Self-notifications are ignored inside CreateNotification.
	if ownerID, err := db.GetPostOwnerID(postID); err == nil {
		_ = db.CreateNotification(ownerID, userID, postID, "post_commented")
	}

	// Redirect back to the post page after successful comment.
	http.Redirect(w, r, "/post?id="+strconv.Itoa(postID), http.StatusSeeOther)
}