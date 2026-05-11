package handlers

import (
	"bytes"
	"database/sql"
	"html/template"
	"net/http"
	"strconv"

	"forum/db"
	"forum/models"
	"forum/sessions"
)

// CommentView extends the Comment model with reaction totals.
// This keeps reaction aggregation out of templates.
type CommentView struct {
	Comment  models.Comment
	Likes    int
	Dislikes int
	IsOwner  bool
}

// SinglePost renders one post page, including:
// - post details + categories
// - comments list
// - like/dislike totals for the post and each comment
func SinglePost(w http.ResponseWriter, r *http.Request) {
	// Validate and parse post ID from query string (?id=123)
	postID, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || postID <= 0 {
		RenderError(w, r, http.StatusBadRequest, "Invalid post ID.")
		return
	}

	// Login state controls whether reaction/comment forms are shown
	currentUserID, loggedIn := sessions.GetUserID(r)

	// AUDIT: moderation visibility rules for single post view:
	// - approved posts are public
	// - owners may see their own posts
	// - moderators/admins may inspect all posts
	isPrivileged := false
	if loggedIn {
		if ok, err := db.IsModeratorOrAdmin(currentUserID); err == nil {
			isPrivileged = ok
		}
	}

	// Load the post and enforce visibility rules
	post, err := db.GetVisiblePostByID(postID, currentUserID, isPrivileged)
	if err != nil {
		if err == sql.ErrNoRows {
			RenderError(w, r, http.StatusNotFound, "Post not found.")
			return
		}
		RenderError(w, r, http.StatusNotFound, "Post not found.")
		return
	}
	post.CreatedAt = FormatDisplayTime(post.CreatedAt)

	// Load all comments for this post (includes username via join)
	comments, err := db.GetCommentsByPostID(postID)
	if err != nil {
		RenderError(w, r, http.StatusInternalServerError, "Could not load comments.")
		return
	}

	// Load reaction totals for the post
	postCounts, err := db.GetPostReactionCounts(postID)
	if err != nil {
		RenderError(w, r, http.StatusInternalServerError, "Could not load post reactions.")
		return
	}

	// Load reaction totals for all comments in one query (avoids N+1 queries)
	commentCounts, err := db.GetCommentCountsByPost(postID)
	if err != nil {
		RenderError(w, r, http.StatusInternalServerError, "Could not load comment reactions.")
		return
	}

	// Build comment view models with formatted timestamps and reaction totals
	views := make([]CommentView, 0, len(comments))
	for _, c := range comments {
		c.CreatedAt = FormatDisplayTime(c.CreatedAt)

		cc := commentCounts[c.ID]
		views = append(views, CommentView{
			Comment:  c,
			Likes:    cc.Likes,
			Dislikes: cc.Dislikes,
			IsOwner:  loggedIn && c.UserID == currentUserID,
		})
	}

	// AUDIT: reporting is available only to privileged roles and only when
	// the post is not already rejected.
	canReport := isPrivileged && post.Status != "rejected"

	// Data passed to the single_post template
	pageData := map[string]interface{}{
		"Post":         post,
		"CommentViews": views,
		"LoggedIn":     loggedIn,
		"PostLikes":    postCounts.Likes,
		"PostDislikes": postCounts.Dislikes,
		"IsPostOwner":  loggedIn && post.UserID == currentUserID,
		"IsPrivileged": isPrivileged,
		"CanReport":    canReport,
	}

	// Render the inner page into a buffer, then wrap it with the global layout
	var buf bytes.Buffer
	if err := Templates.ExecuteTemplate(&buf, "single_post", pageData); err != nil {
		RenderError(w, r, http.StatusInternalServerError, "Template render error.")
		return
	}

	RenderPage(w, r, post.Title, template.HTML(buf.String()))
}