package handlers

import (
	"bytes"
	"html/template"
	"net/http"

	"forum/db"
	"forum/models"
	"forum/sessions"
)

// ActivityReactionView formats reaction activity for templates.
type ActivityReactionView struct {
	PostID      int
	PostTitle   string
	ReactionText string
	CreatedAt   string
}

// ActivityCommentView formats comment activity for templates.
type ActivityCommentView struct {
	PostID     int
	PostTitle  string
	Content    string
	CreatedAt  string
}

// Activity renders the logged-in user's activity page.
// It shows created posts, post reactions and comments.
func Activity(w http.ResponseWriter, r *http.Request) {
	userID, ok := sessions.GetUserID(r)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	posts, err := db.GetPostsByUser(userID)
	if err != nil {
		RenderError(w, r, http.StatusInternalServerError, "Could not load your posts.")
		return
	}

	reactions, err := db.GetReactionsByUser(userID)
	if err != nil {
		RenderError(w, r, http.StatusInternalServerError, "Could not load your reactions.")
		return
	}

	comments, err := db.GetCommentsByUser(userID)
	if err != nil {
		RenderError(w, r, http.StatusInternalServerError, "Could not load your comments.")
		return
	}

	// Format post timestamps for display.
	formattedPosts := make([]models.Post, 0, len(posts))
	for _, p := range posts {
		p.CreatedAt = FormatDisplayTime(p.CreatedAt)
		formattedPosts = append(formattedPosts, p)
	}

	reactionViews := make([]ActivityReactionView, 0, len(reactions))
	for _, r := range reactions {
		txt := "disliked"
		if r.Reaction == 1 {
			txt = "liked"
		}

		reactionViews = append(reactionViews, ActivityReactionView{
			PostID:       r.PostID,
			PostTitle:    r.PostTitle,
			ReactionText: txt,
			CreatedAt:    FormatDisplayTime(r.CreatedAt),
		})
	}

	commentViews := make([]ActivityCommentView, 0, len(comments))
	for _, c := range comments {
		commentViews = append(commentViews, ActivityCommentView{
			PostID:    c.PostID,
			PostTitle: c.PostTitle,
			Content:   c.Content,
			CreatedAt: FormatDisplayTime(c.CreatedAt),
		})
	}

	pageData := map[string]interface{}{
		"Posts":     formattedPosts,
		"Reactions": reactionViews,
		"Comments":  commentViews,
	}

	var buf bytes.Buffer
	if err := Templates.ExecuteTemplate(&buf, "activity", pageData); err != nil {
		RenderError(w, r, http.StatusInternalServerError, "Template render error.")
		return
	}

	RenderPage(w, r, "My Activity", template.HTML(buf.String()))
}