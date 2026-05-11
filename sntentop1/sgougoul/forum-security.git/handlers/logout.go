package handlers

import (
	"net/http"

	"forum/sessions"
)

// Logout handles user logout requests.
// It deletes the session from the database and clears the session cookie.
func Logout(w http.ResponseWriter, r *http.Request) {
	// Allow GET (simple logout link) and POST (more strict setups).
	// Reject any other HTTP method.
	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		RenderError(w, r, http.StatusMethodNotAllowed, "Method not allowed.")
		return
	}

	// Destroy the user's session (removes DB session + cookie).
	sessions.DestroySession(w, r)

	// Redirect the user back to the home page after logout.
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
