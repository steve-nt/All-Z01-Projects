package main

import (
	"log"
	"net/http"

	"forum/db"
	"forum/handlers"
)

func main() {
	// Init DB (creates tables, opens connection)
	db.Init()

	// Init templates
	if err := handlers.InitTemplates(); err != nil {
		log.Fatal("Error parsing templates:", err)
	}

	mux := http.NewServeMux()

	// Static files
	mux.Handle("/static/",
		http.StripPrefix("/static/",
			http.FileServer(http.Dir("static"))))

	// Routes
	mux.HandleFunc("/", handlers.Home)
	mux.HandleFunc("/posts", handlers.Posts)
	mux.HandleFunc("/post", handlers.SinglePost)

	mux.HandleFunc("/register", handlers.Register)
	mux.HandleFunc("/login", handlers.Login)
	mux.HandleFunc("/logout", handlers.Logout)

	mux.HandleFunc("/create-post", handlers.CreatePost)
	mux.HandleFunc("/comment", handlers.CreateComment)

	mux.HandleFunc("/react-post", handlers.ReactPost)
	mux.HandleFunc("/react-comment", handlers.ReactComment)

	// Advanced features
	mux.HandleFunc("/activity", handlers.Activity)
	mux.HandleFunc("/notifications", handlers.Notifications)
	mux.HandleFunc("/edit-post", handlers.EditPost)
	mux.HandleFunc("/delete-post", handlers.DeletePost)
	mux.HandleFunc("/edit-comment", handlers.EditComment)
	mux.HandleFunc("/delete-comment", handlers.DeleteComment)

	// OAuth routes (Google + GitHub)
	mux.HandleFunc("/auth/google", handlers.OAuthGoogleStart)
	mux.HandleFunc("/auth/google/callback", handlers.OAuthGoogleCallback)
	mux.HandleFunc("/auth/github", handlers.OAuthGitHubStart)
	mux.HandleFunc("/auth/github/callback", handlers.OAuthGitHubCallback)

	// AUDIT: moderation routes for moderators/admins.
	mux.HandleFunc("/moderation", handlers.Moderation)
	mux.HandleFunc("/moderation/approve", handlers.ApprovePost)
	mux.HandleFunc("/moderation/reject", handlers.RejectPost)

	// AUDIT: normal user request flow for moderator promotion.
	mux.HandleFunc("/request-moderator", handlers.RequestModerator)

	// AUDIT: moderator/admin report escalation to administrators.
	mux.HandleFunc("/reports/create", handlers.CreateReport)

	// AUDIT: administrator dashboard and management routes.
	mux.HandleFunc("/admin", handlers.AdminDashboard)
	mux.HandleFunc("/admin/update-role", handlers.UpdateUserRole)
	mux.HandleFunc("/admin/requests/approve", handlers.ApproveModeratorRequest)
	mux.HandleFunc("/admin/requests/reject", handlers.RejectModeratorRequest)
	mux.HandleFunc("/admin/categories/create", handlers.CreateCategory)
	mux.HandleFunc("/admin/categories/delete", handlers.DeleteCategory)
	mux.HandleFunc("/admin/reports/resolve", handlers.ResolveReport)

	// Wrap mux with NotFound enforcement
	handler := withNotFound(mux)

	log.Println("Server running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}

func withNotFound(mux *http.ServeMux) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			mux.ServeHTTP(w, r)
			return
		}

		_, pattern := mux.Handler(r)
		if pattern == "/" {
			handlers.RenderError(w, r, http.StatusNotFound, "Page not found.")
			return
		}

		mux.ServeHTTP(w, r)
	})
}