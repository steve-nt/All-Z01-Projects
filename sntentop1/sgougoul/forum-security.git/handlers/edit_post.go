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

// EditPost handles both displaying the edit form and saving changes.
// Only the owner of the post may edit it.
func EditPost(w http.ResponseWriter, r *http.Request) {
	userID, ok := sessions.GetUserID(r)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	postID, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || postID <= 0 {
		RenderError(w, r, http.StatusBadRequest, "Invalid post.")
		return
	}

	post, err := db.GetPostByIDWithCategories(postID)
	if err != nil {
		RenderError(w, r, http.StatusNotFound, "Post not found.")
		return
	}

	// Ownership check
	if post.UserID != userID {
		RenderError(w, r, http.StatusForbidden, "You cannot edit this post.")
		return
	}

	switch r.Method {
	case http.MethodGet:
		pageData := map[string]interface{}{
			"Post": post,
		}

		var buf bytes.Buffer
		if err := Templates.ExecuteTemplate(&buf, "edit_post", pageData); err != nil {
			RenderError(w, r, http.StatusInternalServerError, "Template render error.")
			return
		}

		RenderPage(w, r, "Edit Post", template.HTML(buf.String()))
		return

	case http.MethodPost:
		_ = r.ParseForm()

		title := strings.TrimSpace(r.FormValue("title"))
		content := strings.TrimSpace(r.FormValue("content"))

		if title == "" || content == "" {
			RenderError(w, r, http.StatusBadRequest, "Title and content are required.")
			return
		}

		if err := db.UpdatePost(postID, title, content); err != nil {
			RenderError(w, r, http.StatusInternalServerError, "Could not update post.")
			return
		}

		// AUDIT: after moderation was introduced, edited content should return
		// to the review queue so moderators/admins can re-check the new version.
		if err := db.ResetPostToPending(postID); err != nil {
			RenderError(w, r, http.StatusInternalServerError, "Could not update moderation status.")
			return
		}

		http.Redirect(w, r, "/activity", http.StatusSeeOther)
		return

	default:
		RenderError(w, r, http.StatusMethodNotAllowed, "Method not allowed.")
		return
	}
}

// DeletePost removes a post owned by the current user.
// It also removes related comments, reactions and notifications through the DB layer.
func DeletePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		RenderError(w, r, http.StatusMethodNotAllowed, "Method not allowed.")
		return
	}

	userID, ok := sessions.GetUserID(r)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	postID, err := strconv.Atoi(r.FormValue("post_id"))
	if err != nil || postID <= 0 {
		RenderError(w, r, http.StatusBadRequest, "Invalid post.")
		return
	}

	post, err := db.GetPostByID(postID)
	if err != nil {
		RenderError(w, r, http.StatusNotFound, "Post not found.")
		return
	}

	// AUDIT: moderators and admins are allowed to delete posts as part of forum moderation.
	isPrivileged := false
	if ok, err := db.IsModeratorOrAdmin(userID); err == nil {
		isPrivileged = ok
	}

	if post.UserID != userID && !isPrivileged {
		RenderError(w, r, http.StatusForbidden, "You cannot delete this post.")
		return
	}

	if err := db.DeletePost(postID); err != nil {
		RenderError(w, r, http.StatusInternalServerError, "Could not delete post.")
		return
	}

	http.Redirect(w, r, "/posts", http.StatusSeeOther)
}