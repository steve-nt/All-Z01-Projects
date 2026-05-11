package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func SetupRoutes() {
	// Authentication routes (POST only - GET handled by SPA)
	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			LoginHandler(w, r)
		} else {
			// GET requests handled by SPA
			http.ServeFile(w, r, "frontend/index.html")
		}
	})
	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			RegisterHandler(w, r)
		} else {
			// GET requests handled by SPA
			http.ServeFile(w, r, "frontend/index.html")
		}
	})
	wrapHandler("/logout", LogoutHandler)

	// Post routes
	http.HandleFunc("/new-post", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			CreatePostHandler(w, r)
		} else {
			// GET requests handled by SPA
			http.ServeFile(w, r, "frontend/index.html")
		}
	})
	wrapHandler("/api/posts", PostsAPIHandler)
	wrapHandler("/api/posts/edit", EditPostHandler)
	wrapHandler("/api/posts/delete", DeletePostHandler)

	// Single post view (handled by SPA)
	// API for single post, Since Go doesn't support URL parameters like /api/post/{id} natively
	wrapHandler("/api/post", SinglePostAPIHandler)

	// Comment routes
	wrapHandler("/api/comments/create", CreateCommentHandler)
	wrapHandler("/api/comments", CommentsAPIHandler)
	wrapHandler("/api/comments/edit", EditCommentHandler)
	wrapHandler("/api/comments/delete", DeleteCommentHandler)

	// Like/Dislike routes
	wrapHandler("/api/posts/like", LikePostHandler)
	wrapHandler("/api/comments/like", LikeCommentHandler)

	// User data routes
	http.HandleFunc("/api/user/posts", UserPostsHandler)
	wrapHandler("/api/user/comments", UserCommentsHandler)
	wrapHandler("/api/user/likes", UserLikesHandler)
	wrapHandler("/api/user/dislikes", UserDislikesHandler)
	
	// Category routes
	wrapHandler("/api/categories", CategoriesAPIHandler)

	// Filter routes (for authenticated users)
	wrapHandler("/api/posts/filtered", FilteredPostsHandler)

	// Auth status check - ONLY ONE REGISTRATION
	wrapHandler("/api/auth/status", AuthStatusHandler)

	// Notifications API
	wrapHandler("/api/notifications", NotificationsAPIHandler)
	wrapHandler("/api/notifications/mark-read", MarkNotificationReadHandler)
	wrapHandler("/api/notifications/mark-all-read", MarkAllNotificationsReadHandler)
	wrapHandler("/api/notifications/count", NotificationCountHandler)

	// Profile routes
	wrapHandler("/api/user/profile", ProfileAPIHandler)
	wrapHandler("/profile", ProfileHandler)
	wrapHandler("/profile.html", ProfilePageHandler)

	// Notifications routes (HTML)
	wrapHandler("/notifications", NotificationsPageHandler)
	wrapHandler("/notifications.html", NotificationsPageHandler)

	// Private Messages API routes
	wrapHandler("/api/messages/send", SendMessageHandler)
	wrapHandler("/api/messages", GetMessagesHandler)
	wrapHandler("/api/messages/users", GetUsersListHandler)

	// WebSocket route (real-time)
	// Από εδώ θα συνδέεται το frontend: new WebSocket("ws://localhost:8081/ws")
	http.HandleFunc("/ws", WebSocketHandler)

	// Static files (CSS, images, JavaScript)
	fs := http.FileServer(http.Dir("frontend/"))
	http.Handle("/frontend/", http.StripPrefix("/frontend/", fs))

	// Error routes (for API errors)
	wrapHandler("/400", BadRequestHandler)
	wrapHandler("/404", NotFoundHandler)
	wrapHandler("/500", InternalServerErrorHandler)

	// Catch-all for SPA routes (must be last)
	http.HandleFunc("/", HomeHandler)
}

// wrapHandler wraps handlers with error handling
func wrapHandler(path string, handler http.HandlerFunc) {
	http.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				fmt.Printf("Panic on %s: %v\n", path, err)
				// For API routes, return JSON error
				if strings.HasPrefix(r.URL.Path, "/api/") {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusInternalServerError)
					json.NewEncoder(w).Encode(map[string]interface{}{
						"error":   true,
						"message": "Internal server error",
						"status":  500,
					})
				} else {
					InternalServerErrorHandler(w, r)
				}
			}
		}()
		handler(w, r)
	})
}
