// Package authentication wires HTTP routes for auth, profile, and follow APIs.
// Endpoints: /register, /login, /logout, /api/auth/status, /api/upload-avatar,
// /api/delete-avatar, /api/profile, /api/profile/privacy, /api/follow/request,
// /api/follow/accept, /api/follow/decline, /api/followers, /api/following,
// /api/follow/unfollow, /api/follow/requests, /api/users/search, /api/upload-image.
package authentication

import (
	"net/http"
	"social-network/backend/middleware"
	"social-network/backend/pkg/profile"
)

// SetupAuthRoutes registers all authentication-related routes
// This function should be called from main.go
// All routes use middleware.WrapHandler for logging and error handling
func SetupAuthRoutes() {
	// ===== Public Routes =====
	// These routes are accessible without login.
	http.HandleFunc("/register", middleware.WrapHandler(RegisterHandler))
	http.HandleFunc("/login", middleware.WrapHandler(LoginHandler))

	// ===== Protected Routes =====
	// These routes require the user to be logged in.
	// If not authenticated, redirects to login page.
	http.HandleFunc("/logout", middleware.WrapHandler(middleware.RequireAuth(LogoutHandler)))

	// ===== API Routes =====
	// Auth status can be checked by anyone (returns loggedIn: false if not authenticated).
	http.HandleFunc("/api/auth/status", middleware.WrapHandler(AuthStatusHandler))

	// ===== Media Uploads =====
	// These return JSON errors if not authenticated.
	http.HandleFunc("/api/upload-avatar", middleware.WrapHandler(middleware.RequireAuthJSON(AvatarUploadHandler)))
	http.HandleFunc("/api/delete-avatar", middleware.WrapHandler(middleware.RequireAuthJSON(DeleteAvatarHandler)))
	http.HandleFunc("/api/profile", middleware.WrapHandler(profile.ProfileViewHandler))
	// POST/PUT/PATCH /api/profile/update (auth) -> updates profile fields (nickname, about_me, relationship_status, hobbies, date_of_birth)
	http.HandleFunc("/api/profile/update", middleware.WrapHandler(middleware.RequireAuthJSON(profile.ProfileUpdateHandler)))
	// POST /api/profile/privacy (auth) -> updates Users.is_public
	http.HandleFunc("/api/profile/privacy", middleware.WrapHandler(middleware.RequireAuthJSON(profile.ProfilePrivacyHandler)))
	http.HandleFunc("/api/follow/request", middleware.WrapHandler(middleware.RequireAuthJSON(profile.FollowRequestHandler)))
	http.HandleFunc("/api/follow/accept", middleware.WrapHandler(middleware.RequireAuthJSON(profile.FollowAcceptHandler)))
	http.HandleFunc("/api/follow/decline", middleware.WrapHandler(middleware.RequireAuthJSON(profile.FollowDeclineHandler)))
	http.HandleFunc("/api/followers", middleware.WrapHandler(profile.FollowersHandler))
	http.HandleFunc("/api/following", middleware.WrapHandler(profile.FollowingHandler))
	http.HandleFunc("/api/follow/unfollow", middleware.WrapHandler(middleware.RequireAuthJSON(profile.UnfollowHandler)))
	http.HandleFunc("/api/follow/requests", middleware.WrapHandler(middleware.RequireAuthJSON(profile.FollowRequestsHandler)))
	http.HandleFunc("/api/users/search", middleware.WrapHandler(middleware.RequireAuthJSON(profile.UsersSearchHandler)))

	// NOTE: ImageUploadHandler is for Part 3 (Posts & Groups)
	// Currently a placeholder - Part 3 will implement full post image storage
	http.HandleFunc("/api/upload-image", middleware.WrapHandler(middleware.RequireAuthJSON(ImageUploadHandler)))
}
