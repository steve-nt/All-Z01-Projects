package posts

import (
	"net/http"
	"social-network/backend/middleware"
)

// SetupPostRoutes registers Posts & Comments API endpoints.
// Call this from main.go.
func SetupPostRoutes() {
	// Feed + single post fetch (auth optional; privacy filtering applies).
	http.HandleFunc("/api/posts", middleware.WrapHandler(PostsHandler))
	http.HandleFunc("/api/posts/view", middleware.WrapHandler(PostViewHandler))

	// Post creation requires auth.
	http.HandleFunc("/api/posts/create", middleware.WrapHandler(middleware.RequireAuthJSON(CreatePostHandler)))

	// Comments
	http.HandleFunc("/api/posts/comments", middleware.WrapHandler(CommentsHandler)) // GET list, POST create (auth required for create)
}

