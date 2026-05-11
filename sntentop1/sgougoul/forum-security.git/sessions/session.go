package sessions

import (
	"database/sql"
	"net/http"
	"os"
	"strings"
	"time"

	"forum/db"

	"github.com/google/uuid"
)

// cookieName is the name of the browser cookie storing the session token.
const cookieName = "forum_session"

// sessionDuration defines how long a login session stays valid.
const sessionDuration = 24 * time.Hour

// shouldUseSecureCookies determines whether session cookies should be marked Secure.
// AUDIT:
// - HTTPS requests automatically use Secure cookies
// - SESSION_COOKIE_SECURE=true can force Secure cookies when TLS is terminated upstream
func shouldUseSecureCookies(r *http.Request) bool {
	if r != nil && r.TLS != nil {
		return true
	}

	v := strings.ToLower(strings.TrimSpace(os.Getenv("SESSION_COOKIE_SECURE")))
	return v == "1" || v == "true" || v == "yes"
}

// sessionCookie builds the session cookie with consistent security settings.
// AUDIT:
// - HttpOnly prevents JavaScript access to the session token
// - SameSite=Lax helps reduce CSRF risk on top-level navigation
// - Secure is enabled for HTTPS, or can be forced by environment when deployed behind TLS proxy
func sessionCookie(r *http.Request, value string, expiresAt time.Time, maxAge int) *http.Cookie {
	return &http.Cookie{
		Name:     cookieName,
		Value:    value,
		Path:     "/",
		Expires:  expiresAt,
		MaxAge:   maxAge,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   shouldUseSecureCookies(r),
	}
}

// expiredSessionCookie returns a cookie that immediately removes the session
// from the browser using the same policy shape as the active cookie.
func expiredSessionCookie(r *http.Request) *http.Cookie {
	return sessionCookie(r, "", time.Unix(0, 0), -1)
}

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
	// AUDIT: rotate any prior browser cookie first to reduce session fixation risk.
	// The server-side session store remains the source of truth.
	http.SetCookie(w, expiredSessionCookie(r))

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

	// Send cookie to browser
	http.SetCookie(w, sessionCookie(r, token, expiresAt, int(sessionDuration.Seconds())))

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
	http.SetCookie(w, expiredSessionCookie(r))
}

// deleteByToken removes a session row from DB.
func deleteByToken(token string) error {
	_, err := db.DB.Exec(`DELETE FROM sessions WHERE token = ?`, token)
	return err
}