package handlers

import (
	"net/http"

	"forum/db"
	"forum/sessions"
)

// requireLogin ensures the request belongs to an authenticated user.
// Returns (userID, true) when logged in, otherwise redirects to /login.
func requireLogin(w http.ResponseWriter, r *http.Request) (int, bool) {
	userID, ok := sessions.GetUserID(r)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return 0, false
	}
	return userID, true
}

// requireAdmin ensures the authenticated user is an administrator.
// Returns true only when the user is logged in and has admin role.
func requireAdmin(w http.ResponseWriter, r *http.Request) (int, bool) {
	userID, ok := requireLogin(w, r)
	if !ok {
		return 0, false
	}

	isAdmin, err := db.IsAdmin(userID)
	if err != nil {
		RenderError(w, r, http.StatusInternalServerError, "Could not verify admin access.")
		return 0, false
	}
	if !isAdmin {
		RenderError(w, r, http.StatusForbidden, "Only administrators can access this page.")
		return 0, false
	}

	return userID, true
}

// requireModeratorOrAdmin ensures the authenticated user has moderation access.
// Returns true only when the user is logged in and is moderator/admin.
func requireModeratorOrAdmin(w http.ResponseWriter, r *http.Request) (int, bool) {
	userID, ok := requireLogin(w, r)
	if !ok {
		return 0, false
	}

	isPrivileged, err := db.IsModeratorOrAdmin(userID)
	if err != nil {
		RenderError(w, r, http.StatusInternalServerError, "Could not verify access level.")
		return 0, false
	}
	if !isPrivileged {
		RenderError(w, r, http.StatusForbidden, "You do not have access to moderation.")
		return 0, false
	}

	return userID, true
}