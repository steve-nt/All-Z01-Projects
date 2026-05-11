package authentication

import (
	"encoding/json"
	"net/http"
	"social-network/backend/pkg/db/sqlite"
	"social-network/backend/utils"
	"strings"
)

// wantsJSON returns true when the client expects JSON (SPA / API usage).
// We use this to keep HTML template flows working while adding SPA-friendly responses.
func wantsJSON(r *http.Request) bool {
	accept := r.Header.Get("Accept")
	contentType := r.Header.Get("Content-Type")
	xrw := r.Header.Get("X-Requested-With")
	return strings.Contains(accept, "application/json") ||
		strings.Contains(contentType, "application/json") ||
		xrw == "XMLHttpRequest"
}

func writeJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

// LogoutHandler handles user logout
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session")
	if err == nil {
		db := sqlite.GetDB()
		db.Exec("DELETE FROM Sessions WHERE cookie_value = ?", cookie.Value)
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	})

	if wantsJSON(r) {
		writeJSON(w, http.StatusOK, map[string]any{"ok": true})
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// AuthStatusHandler checks authentication status for API calls
func AuthStatusHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session")
	isLoggedIn := err == nil && cookie != nil && utils.IsValidSession(cookie.Value)

	if isLoggedIn {
		userID := utils.GetUserIDFromSession(cookie.Value)
		nickname := utils.GetNicknameFromSession(cookie.Value)
		writeJSON(w, http.StatusOK, map[string]interface{}{
			"loggedIn": true,
			"userID":   userID,
			"nickname": nickname,
		})
	} else {
		writeJSON(w, http.StatusOK, map[string]interface{}{
			"loggedIn": false,
		})
	}
}
