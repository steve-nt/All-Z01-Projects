package server

import (
	"net/http"

	"forum/handlers"
)

// registerRoutes defines all application endpoints.
// AUDIT: sensitive routes are protected with rate limiting.
func registerRoutes(mux *http.ServeMux) {

	mux.HandleFunc("/", handlers.Home)
	mux.HandleFunc("/posts", handlers.Posts)
	mux.HandleFunc("/post", handlers.SinglePost)

	mux.HandleFunc("/register", handlers.RateLimit(handlers.Register, 5, 0.5))
	mux.HandleFunc("/login", handlers.RateLimit(handlers.Login, 5, 0.5))
	mux.HandleFunc("/logout", handlers.Logout)

	mux.HandleFunc("/create-post", handlers.RateLimit(handlers.CreatePost, 10, 1.0))
	mux.HandleFunc("/comment", handlers.RateLimit(handlers.CreateComment, 10, 1.0))

	mux.HandleFunc("/react-post", handlers.RateLimit(handlers.ReactPost, 20, 2.0))
	mux.HandleFunc("/react-comment", handlers.RateLimit(handlers.ReactComment, 20, 2.0))

	mux.HandleFunc("/activity", handlers.Activity)
	mux.HandleFunc("/notifications", handlers.Notifications)
	mux.HandleFunc("/edit-post", handlers.EditPost)
	mux.HandleFunc("/delete-post", handlers.DeletePost)
	mux.HandleFunc("/edit-comment", handlers.EditComment)
	mux.HandleFunc("/delete-comment", handlers.DeleteComment)

	mux.HandleFunc("/auth/google", handlers.RateLimit(handlers.OAuthGoogleStart, 5, 0.5))
	mux.HandleFunc("/auth/google/callback", handlers.OAuthGoogleCallback)
	mux.HandleFunc("/auth/github", handlers.RateLimit(handlers.OAuthGitHubStart, 5, 0.5))
	mux.HandleFunc("/auth/github/callback", handlers.OAuthGitHubCallback)

	mux.HandleFunc("/moderation", handlers.Moderation)
	mux.HandleFunc("/moderation/approve", handlers.RateLimit(handlers.ApprovePost, 10, 1.0))
	mux.HandleFunc("/moderation/reject", handlers.RateLimit(handlers.RejectPost, 10, 1.0))

	mux.HandleFunc("/request-moderator", handlers.RateLimit(handlers.RequestModerator, 3, 0.2))
	mux.HandleFunc("/reports/create", handlers.RateLimit(handlers.CreateReport, 10, 1.0))

	mux.HandleFunc("/admin", handlers.AdminDashboard)
	mux.HandleFunc("/admin/update-role", handlers.RateLimit(handlers.UpdateUserRole, 10, 1.0))
	mux.HandleFunc("/admin/requests/approve", handlers.RateLimit(handlers.ApproveModeratorRequest, 10, 1.0))
	mux.HandleFunc("/admin/requests/reject", handlers.RateLimit(handlers.RejectModeratorRequest, 10, 1.0))
	mux.HandleFunc("/admin/categories/create", handlers.RateLimit(handlers.CreateCategory, 10, 1.0))
	mux.HandleFunc("/admin/categories/delete", handlers.RateLimit(handlers.DeleteCategory, 10, 1.0))
	mux.HandleFunc("/admin/reports/resolve", handlers.RateLimit(handlers.ResolveReport, 10, 1.0))
}