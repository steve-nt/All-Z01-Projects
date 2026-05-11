package handlers

import (
	"bytes"
	"html/template"
	"net/http"
	"strings"

	"forum/db"
	"forum/sessions"
)

// HomePostView contains the fields needed by the home page template.
type HomePostView struct {
	ID         int
	Title      string
	Username   string
	Preview    string
	CreatedAt  string
	Categories []string
	Likes      int
	Dislikes   int
}

// Home renders the forum landing page.
// It shows a welcome section, forum guidance and the top liked posts.
func Home(w http.ResponseWriter, r *http.Request) {
	_, loggedIn := sessions.GetUserID(r)

	// Load the top 5 most liked posts for the home page.
	topPosts, err := db.GetTopLikedPosts(5)
	if err != nil {
		RenderError(w, r, http.StatusInternalServerError, "Failed to load posts.")
		return
	}

	// Load reaction totals for the selected posts in one query.
	postIDs := make([]int, 0, len(topPosts))
	for _, p := range topPosts {
		postIDs = append(postIDs, p.ID)
	}

	countsMap, err := db.GetPostCountsByPostIDs(postIDs)
	if err != nil {
		RenderError(w, r, http.StatusInternalServerError, "Failed to load reactions.")
		return
	}

	views := make([]HomePostView, 0, len(topPosts))
	for _, p := range topPosts {
		preview := strings.TrimSpace(p.Content)
		if len(preview) > 160 {
			preview = preview[:160] + "..."
		}

		counts := countsMap[p.ID]

		views = append(views, HomePostView{
			ID:         p.ID,
			Title:      p.Title,
			Username:   p.Username,
			Preview:    preview,
			CreatedAt:  FormatDisplayTime(p.CreatedAt),
			Categories: p.Categories,
			Likes:      counts.Likes,
			Dislikes:   counts.Dislikes,
		})
	}

	pageData := map[string]interface{}{
		"LoggedIn": loggedIn,
		"TopPosts": views,
	}

	var buf bytes.Buffer
	if err := Templates.ExecuteTemplate(&buf, "home", pageData); err != nil {
		RenderError(w, r, http.StatusInternalServerError, "Template render error.")
		return
	}

	RenderPage(w, r, "Forum Home", template.HTML(buf.String()))
}