package handlers

import (
	"bytes"
	"html/template"
	"net/http"
	"strings"

	"forum/db"
	"forum/models"
	"forum/sessions"
)

// PostView extends the DB Post model with reaction counts.
// This keeps database logic separate from display logic.
type PostView struct {
	models.Post
	Likes    int
	Dislikes int
}

// Posts renders the main posts listing page.
// Supports filtering by category, user's own posts, or liked posts.
func Posts(w http.ResponseWriter, r *http.Request) {
	// Detect login status (needed for protected filters)
	userID, loggedIn := sessions.GetUserID(r)

	// Read optional query parameters
	category := strings.TrimSpace(r.URL.Query().Get("category"))
	filter := strings.ToLower(strings.TrimSpace(r.URL.Query().Get("filter")))

	// Load all categories for the filter UI
	cats, err := db.GetAllCategories()
	if err != nil {
		RenderError(w, r, http.StatusInternalServerError, "Could not load categories.")
		return
	}

	// Load posts depending on selected filter
	var postsList []models.Post

	switch filter {

	case "mine":
		// Personal filters require login
		if !loggedIn {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		postsList, err = db.GetPostsByUser(userID)

	case "liked":
		// Personal filters require login
		if !loggedIn {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		postsList, err = db.GetLikedPostsByUser(userID)

	default:
		// Default view = publicly visible posts only
		postsList, err = db.GetAllPosts(category)
	}

	if err != nil {
		RenderError(w, r, http.StatusInternalServerError, "Failed to load posts.")
		return
	}

	// Collect all post IDs so reaction counts can be fetched in ONE query
	// (prevents N+1 query problem and improves performance)
	postIDs := make([]int, 0, len(postsList))
	for _, p := range postsList {
		postIDs = append(postIDs, p.ID)
	}

	countsMap, err := db.GetPostCountsByPostIDs(postIDs)
	if err != nil {
		RenderError(w, r, http.StatusInternalServerError, "Failed to load reactions.")
		return
	}

	// Build view objects with formatted timestamps and reaction totals
	views := make([]PostView, 0, len(postsList))
	for _, p := range postsList {
		c := countsMap[p.ID]

		p.CreatedAt = FormatDisplayTime(p.CreatedAt)

		views = append(views, PostView{
			Post:     p,
			Likes:    c.Likes,
			Dislikes: c.Dislikes,
		})
	}

	// Data passed to the template
	pageData := map[string]interface{}{
		"LoggedIn":   loggedIn,
		"Category":   category,
		"Filter":     filter,
		"Categories": cats,
		"Posts":      views,
	}

	// Render page content into a buffer first
	// (allows safe wrapping inside the global layout)
	var buf bytes.Buffer
	if err := Templates.ExecuteTemplate(&buf, "posts", pageData); err != nil {
		RenderError(w, r, http.StatusInternalServerError, "Template render error.")
		return
	}

	// Inject rendered content into layout template
	RenderPage(w, r, "Posts", template.HTML(buf.String()))
}