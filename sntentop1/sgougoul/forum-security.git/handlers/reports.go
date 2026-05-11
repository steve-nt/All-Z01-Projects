package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"forum/db"
	"forum/sessions"
)

// CreateReport lets a moderator report a post to administrators.
func CreateReport(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		RenderError(w, r, http.StatusMethodNotAllowed, "Method not allowed.")
		return
	}

	userID, ok := sessions.GetUserID(r)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	role, err := db.GetUserRole(userID)
	if err != nil {
		RenderError(w, r, http.StatusInternalServerError, "Could not verify role.")
		return
	}
	if role != "moderator" && role != "admin" {
		RenderError(w, r, http.StatusForbidden, "Only moderators and admins can report posts.")
		return
	}

	postID, err := strconv.Atoi(r.FormValue("post_id"))
	if err != nil || postID <= 0 {
		RenderError(w, r, http.StatusBadRequest, "Invalid post.")
		return
	}

	reason := strings.TrimSpace(r.FormValue("reason"))
	if reason == "" {
		RenderError(w, r, http.StatusBadRequest, "Report reason is required.")
		return
	}

	// AUDIT: avoid repeated pending escalation for the same unresolved post.
	hasPending, err := db.HasPendingReportForPost(postID)
	if err != nil {
		RenderError(w, r, http.StatusInternalServerError, "Could not validate report state.")
		return
	}
	if hasPending {
		RenderError(w, r, http.StatusBadRequest, "This post already has a pending report.")
		return
	}

	if err := db.CreateReport(postID, userID, reason); err != nil {
		RenderError(w, r, http.StatusInternalServerError, "Could not create report.")
		return
	}

	http.Redirect(w, r, "/post?id="+strconv.Itoa(postID), http.StatusSeeOther)
}