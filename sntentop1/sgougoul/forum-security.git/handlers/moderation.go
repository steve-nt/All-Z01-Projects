package handlers

import (
	"bytes"
	"html/template"
	"net/http"
	"strconv"

	"forum/db"
)

// Moderation renders the moderation queue for moderators and admins.
// Only privileged roles may access this page.
func Moderation(w http.ResponseWriter, r *http.Request) {
	_, ok := requireModeratorOrAdmin(w, r)
	if !ok {
		return
	}

	posts, err := db.GetPendingPosts()
	if err != nil {
		RenderError(w, r, http.StatusInternalServerError, "Could not load moderation queue.")
		return
	}

	for i := range posts {
		posts[i].CreatedAt = FormatDisplayTime(posts[i].CreatedAt)
	}

	pageData := map[string]interface{}{
		"Posts": posts,
	}

	var buf bytes.Buffer
	if err := Templates.ExecuteTemplate(&buf, "moderation", pageData); err != nil {
		RenderError(w, r, http.StatusInternalServerError, "Template render error.")
		return
	}

	RenderPage(w, r, "Moderation Queue", template.HTML(buf.String()))
}

// ApprovePost handles approval of a pending post.
// Only moderators/admins may approve content.
func ApprovePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		RenderError(w, r, http.StatusMethodNotAllowed, "Method not allowed.")
		return
	}

	userID, ok := requireModeratorOrAdmin(w, r)
	if !ok {
		return
	}

	postID, err := strconv.Atoi(r.FormValue("post_id"))
	if err != nil || postID <= 0 {
		RenderError(w, r, http.StatusBadRequest, "Invalid post.")
		return
	}

	if err := db.ApprovePost(postID, userID); err != nil {
		RenderError(w, r, http.StatusInternalServerError, "Could not approve post.")
		return
	}

	http.Redirect(w, r, "/moderation", http.StatusSeeOther)
}

// RejectPost handles rejection of a pending post.
// Only moderators/admins may reject content.
func RejectPost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		RenderError(w, r, http.StatusMethodNotAllowed, "Method not allowed.")
		return
	}

	userID, ok := requireModeratorOrAdmin(w, r)
	if !ok {
		return
	}

	postID, err := strconv.Atoi(r.FormValue("post_id"))
	if err != nil || postID <= 0 {
		RenderError(w, r, http.StatusBadRequest, "Invalid post.")
		return
	}

	if err := db.RejectPost(postID, userID); err != nil {
		RenderError(w, r, http.StatusInternalServerError, "Could not reject post.")
		return
	}

	http.Redirect(w, r, "/moderation", http.StatusSeeOther)
}