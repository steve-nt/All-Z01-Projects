package sessions

import (
	"database/sql"
	"net/http"
	"time"

	"forum/db"

	"github.com/google/uuid"
)

// cookieName is the name of the browser cookie storing the session token.
const cookieName = "forum_session"

// sessionDuration defines how long a login session stays valid.
const sessionDuration = 24 * time.Hour

// CreateSession creates a new login session for the given user.
//
// IMPORTANT FEATURE:
// Only ONE active session per user is allowed.
// If the user logs in from another browser,
// the previous session is deleted.
//
// Steps:
//  1. delete existing session for that user
//  2. insert a new session row
//  3. set a cookie containing the session token
func CreateSession(w http.ResponseWriter, r *http.Request, userID int) error {
	// Generate a random unique session token
	token := uuid.NewString()

	// Calculate expiration timestamp
	expiresAt := time.Now().Add(sessionDuration)

	// Use a transaction so delete+insert happens atomically
	tx, err := db.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Remove any existing session for this user
	// (sessions.user_id column is UNIQUE)
	if _, err := tx.Exec(`DELETE FROM sessions WHERE user_id = ?`, userID); err != nil {
		return err
	}

	// Insert the new session row
	if _, err := tx.Exec(
		`INSERT INTO sessions (user_id, token, expires_at) VALUES (?, ?, ?)`,
		userID, token, expiresAt,
	); err != nil {
		return err
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return err
	}

	// Determine if cookie should be Secure (HTTPS only)
	secure := false
	if r != nil && r.TLS != nil {
		secure = true
	}

	// Send cookie to browser
	http.SetCookie(w, &http.Cookie{
		Name:     cookieName,
		Value:    token,
		Path:     "/",
		Expires:  expiresAt,
		MaxAge:   int(sessionDuration.Seconds()),
		HttpOnly: true, // prevents JS access (XSS protection)
		SameSite: http.SameSiteLaxMode,
		Secure:   secure, // only true when using HTTPS
	})

	return nil
}

// GetUserID checks the request cookie and returns:
//
//	(userID, true)  -> valid logged-in user
//	(0, false)      -> not logged in or session expired
//
// It validates:
//   - cookie exists
//   - session exists in DB
//   - session not expired
func GetUserID(r *http.Request) (int, bool) {
	// Read session cookie
	c, err := r.Cookie(cookieName)
	if err != nil || c.Value == "" {
		return 0, false
	}

	var userID int
	var expiresAt time.Time

	// Lookup session in database
	err = db.DB.QueryRow(
		`SELECT user_id, expires_at FROM sessions WHERE token = ?`,
		c.Value,
	).Scan(&userID, &expiresAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, false
		}
		return 0, false
	}

	// If expired → delete session and treat as logged out
	if time.Now().After(expiresAt) {
		_ = deleteByToken(c.Value)
		return 0, false
	}

	return userID, true
}

// DestroySession logs the user out.
//
// It:
//  1. deletes the session row from DB
//  2. clears the browser cookie
func DestroySession(w http.ResponseWriter, r *http.Request) {
	// If cookie exists, delete DB session
	c, err := r.Cookie(cookieName)
	if err == nil && c.Value != "" {
		_ = deleteByToken(c.Value)
	}

	// Clear cookie in browser
	http.SetCookie(w, &http.Cookie{
		Name:     cookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1, // delete immediately
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
}

// deleteByToken removes a session row from DB.
func deleteByToken(token string) error {
	_, err := db.DB.Exec(`DELETE FROM sessions WHERE token = ?`, token)
	return err
}
