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

// CreatePost handles both displaying the post creation page (GET)
// and inserting a new post into the database (POST).
// Only logged-in users are allowed to create posts.
func CreatePost(w http.ResponseWriter, r *http.Request) {
	// Verify user session.
	userID, ok := sessions.GetUserID(r)
	if !ok {

		// Anonymous users must login first.
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	switch r.Method {

	// ----------------------------
	// GET → Show post creation form
	// ----------------------------
	case http.MethodGet:

		// Load premade categories from DB.
		cats, err := db.GetAllCategories()
		if err != nil {
			RenderError(w, r, http.StatusInternalServerError, "Could not load categories.")
			return
		}

		// Pass categories to template.
		pageData := map[string]interface{}{
			"Categories": cats,
		}

		// Render template into buffer first.
		var buf bytes.Buffer
		if err := Templates.ExecuteTemplate(&buf, "create_post", pageData); err != nil {
			RenderError(w, r, http.StatusInternalServerError, "Template render error.")
			return
		}

		// Inject into layout wrapper.
		RenderPage(w, r, "Create Post", template.HTML(buf.String()))
		return

	// ----------------------------
	// POST → Create the post
	// ----------------------------
	case http.MethodPost:

		_ = r.ParseForm()

		// Trim whitespace from inputs.
		title := strings.TrimSpace(r.FormValue("title"))
		content := strings.TrimSpace(r.FormValue("content"))

		// Prevent empty posts.
		if title == "" || content == "" {
			RenderError(w, r, http.StatusBadRequest, "Title and content are required.")
			return
		}

		// Read selected category IDs from checkboxes.
		// HTML uses: name="category_ids"
		var categoryIDs []int

		for _, s := range r.Form["category_ids"] {

			id, err := strconv.Atoi(s)
			if err != nil {
				continue
			}

			if id > 0 {
				categoryIDs = append(categoryIDs, id)
			}
		}

		// Enforce at least one category (project requirement).
		if len(categoryIDs) == 0 {
			RenderError(w, r, http.StatusBadRequest, "Please select at least one category.")
			return
		}

		// Insert post into DB.
		// AUDIT: new posts are explicitly created as pending in the DB layer,
		// so they must be reviewed before becoming publicly visible.
		postID, err := db.CreatePost(userID, title, content)
		if err != nil {
			RenderError(w, r, http.StatusInternalServerError, "Could not create post.")
			return
		}

		// Link post to selected premade categories.
		if err := db.SetPostCategoryIDs(postID, categoryIDs); err != nil {
			RenderError(w, r, http.StatusInternalServerError, "Could not save categories.")
			return
		}

		// AUDIT: notify every moderator/admin that a post is waiting for approval.
		reviewerIDs, err := db.GetModeratorAndAdminUserIDs()
		if err == nil {
			for _, reviewerID := range reviewerIDs {
				_ = db.CreateCustomNotification(
					reviewerID,
					userID,
					postID,
					"post_pending_review",
					"A new post is waiting for moderation review.",
				)
			}
		}

		// Redirect to the user's activity page so they can still find
		// their newly submitted pending post.
		http.Redirect(w, r, "/activity", http.StatusSeeOther)
		return

	default:

		RenderError(w, r, http.StatusMethodNotAllowed, "Method not allowed.")
		return
	}
}