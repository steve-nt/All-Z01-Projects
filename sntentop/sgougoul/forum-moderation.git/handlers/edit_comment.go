package handlers

import (
	"bytes"
	"html/template"
	"net/http"
	"strconv"
	"strings"

	"forum/db"
	"forum/sessions"
)

// EditComment displays the edit form and saves changes for a comment.
// Only the comment owner may edit it.
func EditComment(w http.ResponseWriter, r *http.Request) {
	userID, ok := sessions.GetUserID(r)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	commentID, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || commentID <= 0 {
		RenderError(w, r, http.StatusBadRequest, "Invalid comment.")
		return
	}

	comment, err := db.GetCommentByID(commentID)
	if err != nil {
		RenderError(w, r, http.StatusNotFound, "Comment not found.")
		return
	}

	if comment.UserID != userID {
		RenderError(w, r, http.StatusForbidden, "You cannot edit this comment.")
		return
	}

	switch r.Method {
	case http.MethodGet:
		pageData := map[string]interface{}{
			"Comment": comment,
		}

		var buf bytes.Buffer
		if err := Templates.ExecuteTemplate(&buf, "edit_comment", pageData); err != nil {
			RenderError(w, r, http.StatusInternalServerError, "Template render error.")
			return
		}

		RenderPage(w, r, "Edit Comment", template.HTML(buf.String()))
		return

	case http.MethodPost:
		_ = r.ParseForm()

		content := strings.TrimSpace(r.FormValue("content"))
		if content == "" {
			RenderError(w, r, http.StatusBadRequest, "Comment cannot be empty.")
			return
		}

		if err := db.UpdateComment(commentID, content); err != nil {
			RenderError(w, r, http.StatusInternalServerError, "Could not update comment.")
			return
		}

		http.Redirect(w, r, "/post?id="+strconv.Itoa(comment.PostID), http.StatusSeeOther)
		return

	default:
		RenderError(w, r, http.StatusMethodNotAllowed, "Method not allowed.")
		return
	}
}

// DeleteComment removes a comment owned by the current user.
func DeleteComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		RenderError(w, r, http.StatusMethodNotAllowed, "Method not allowed.")
		return
	}

	userID, ok := sessions.GetUserID(r)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	commentID, err := strconv.Atoi(r.FormValue("comment_id"))
	if err != nil || commentID <= 0 {
		RenderError(w, r, http.StatusBadRequest, "Invalid comment.")
		return
	}

	comment, err := db.GetCommentByID(commentID)
	if err != nil {
		RenderError(w, r, http.StatusNotFound, "Comment not found.")
		return
	}

	// AUDIT: moderators and admins are allowed to remove comments during forum moderation.
	isPrivileged := false
	if ok, err := db.IsModeratorOrAdmin(userID); err == nil {
		isPrivileged = ok
	}

	if comment.UserID != userID && !isPrivileged {
		RenderError(w, r, http.StatusForbidden, "You cannot delete this comment.")
		return
	}

	if err := db.DeleteComment(commentID); err != nil {
		RenderError(w, r, http.StatusInternalServerError, "Could not delete comment.")
		return
	}

	http.Redirect(w, r, "/post?id="+strconv.Itoa(comment.PostID), http.StatusSeeOther)
}