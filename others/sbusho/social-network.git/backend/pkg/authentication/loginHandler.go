package authentication

import (
	"fmt"
	"net/http"
	"social-network/backend/pkg/db/sqlite"
	"social-network/backend/utils"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		// Check for success message from registration
		successType := r.URL.Query().Get("success")
		message := r.URL.Query().Get("message")

		data := make(map[string]interface{})

		if successType == "registration" {
			data["SuccessMessage"] = "Registration successful! Please log in with your new account."
		} else if message != "" {
			data["SuccessMessage"] = message
		}

		utils.FileService("login.html", w, data)
		return
	}

	db := sqlite.GetDB()

	email := strings.TrimSpace(r.FormValue("email"))
	password := r.FormValue("password")

	if email == "" || password == "" {
		if wantsJSON(r) {
			writeJSON(w, http.StatusBadRequest, map[string]any{"error": "Email and password cannot be empty"})
			return
		}
		utils.FileService("login.html", w, map[string]interface{}{
			"ErrorMessage": "Email and password cannot be empty",
		})
		return
	}

	var userID int
	var passwordHash string

	err := db.QueryRow(
		"SELECT user_id, password_hash FROM Users WHERE email = ?",
		email,
	).Scan(&userID, &passwordHash)

	if err != nil || bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password)) != nil {
		if wantsJSON(r) {
			writeJSON(w, http.StatusUnauthorized, map[string]any{"error": "Invalid email or password"})
			return
		}
		utils.FileService("login.html", w, map[string]interface{}{
			"ErrorMessage": "Invalid email or password",
		})
		return
	}

	// Cleanup old sessions for this user
	_, err = db.Exec("DELETE FROM Sessions WHERE user_id = ?", userID)
	if err != nil {
		fmt.Printf("Warning: Failed to cleanup old sessions for user %d: %v\n", userID, err)
	}

	// Create secure session cookie
	cookieValue := utils.GenerateCookieValue()
	expiration := time.Now().Add(24 * time.Hour)

	// Set session cookie with SameSite
	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    cookieValue,
		Path:     "/",
		Expires:  expiration,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})

	// Store session in database
	if _, err = sqlite.Insert(db, "Sessions", "(user_id, cookie_value, expiration_date)", userID, cookieValue, expiration); err != nil {
		if wantsJSON(r) {
			writeJSON(w, http.StatusInternalServerError, map[string]any{"error": "Error creating session"})
			return
		}
		http.Error(w, "Error creating session", http.StatusInternalServerError)
		return
	}

	// SPA/API: return JSON instead of redirect
	if wantsJSON(r) {
		writeJSON(w, http.StatusOK, map[string]any{"ok": true})
		return
	}

	// HTML flow: Redirect to home
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
