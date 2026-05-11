// Package utils provides shared helpers for HTTP handling, validation, sessions, and formatting.
/*
Responsibilities:
- Serve the SPA shell for routes the frontend owns.
- Provide validation helpers (email/password), session lookups, and common formatting utilities.
- Offer simple routing helpers (API path detection).
*/

package utils

import (
	"crypto/rand"
	"fmt"
	"io"
	"net/http"
	"os"
	"realtimeforum/internals/database"
	"regexp"
	"strings"
	"time"
)

// TemplateData holds data to pass to templates
type TemplateData struct {
	IsLoggedIn bool
	Username   string
	UserID     int
	Message    string
	Error      string
	Data       interface{}
}

// =========================
// Section: SPA shell serving
// =========================

// FileService always serves the SPA entrypoint regardless of filename.
// Used by routes where the frontend router handles rendering. Ignores the filename parameter intentionally.
// Misuse warning: only use for routes intended to defer to the SPA; does not enforce auth.
func FileService(filename string, w http.ResponseWriter, data any) {
	// For SPA, always serve index.html (JavaScript handles routing)
	file, err := os.Open("frontend/index.html")
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}
	defer file.Close()
	
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	io.Copy(w, file)
}

// FileServiceWithAuth serves the SPA with authentication context
// Same as FileService but preserves request context/cookies; still does not enforce auth on its own.
func FileServiceWithAuth(filename string, w http.ResponseWriter, r *http.Request, data interface{}) {
	// For SPA, always serve index.html (JavaScript handles routing and auth)
	file, err := os.Open("frontend/index.html")
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}
	defer file.Close()
	
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	io.Copy(w, file)
}

// ======================
// Section: Validation
// ======================

// IsValidEmail performs a lightweight validation for email format.
// Returns false for clearly invalid addresses; not a full RFC validator.
// Typical callers: registration, password reset handlers.
func IsValidEmail(email string) bool {
	// Trim whitespace and reject obviously malformed addresses early
	email = strings.TrimSpace(email)

	// Basic checks
	if len(email) < 5 || len(email) > 254 {
		return false
	}

	// Must contain exactly one @
	if strings.Count(email, "@") != 1 {
		return false
	}

	// Split and check parts exist
	parts := strings.Split(email, "@")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return false
	}

	localPart := parts[0]
	domainPart := parts[1]

	// Local part length check
	if len(localPart) > 64 {
		return false
	}

	// Domain must contain at least one dot
	if !strings.Contains(domainPart, ".") {
		return false
	}

	// Domain can't start or end with dot or dash
	if strings.HasPrefix(domainPart, ".") || strings.HasSuffix(domainPart, ".") ||
		strings.HasPrefix(domainPart, "-") || strings.HasSuffix(domainPart, "-") {
		return false
	}

	// Very simple regex - just basic characters
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, err := regexp.MatchString(emailRegex, email)

	return err == nil && matched

}

// TestBasicEmailValidation is a manual helper to print sample validation results.
// Not used by handlers; safe to leave as-is for local diagnostics.
func TestBasicEmailValidation() {
	testEmails := []string{
		"yuki@gmail.com",         // Should be true
		"test@example.com",       // Should be true
		"user.name@domain.co.uk", // Should be true
		"invalid-email",          // Should be false
		"@invalid.com",           // Should be false
		"invalid@",               // Should be false
		"invalid@@domain.com",    // Should be false
		"test@domain",            // Should be false (no TLD)
		"",                       // Should be false
	}

	fmt.Println("=== BASIC EMAIL VALIDATION TEST ===")
	for _, email := range testEmails {
		valid := IsValidEmail(email)
		fmt.Printf("%-25s -> %v\n", email, valid)
	}
}

// IsValidPassword enforces minimum length and character class presence.
// Returns false if any required class (lower, upper, digit, special) is missing.
func IsValidPassword(password string) bool {
	if len(password) < 8 {
		return false
	}

	// Check for at least one lowercase letter
	hasLower, _ := regexp.MatchString(`[a-z]`, password)
	if !hasLower {
		return false
	}

	// Check for at least one uppercase letter
	hasUpper, _ := regexp.MatchString(`[A-Z]`, password)
	if !hasUpper {
		return false
	}

	// Check for at least one digit
	hasDigit, _ := regexp.MatchString(`[0-9]`, password)
	if !hasDigit {
		return false
	}

	// Check for at least one special character
	hasSpecial, _ := regexp.MatchString(`[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?~`+"`"+`]`, password)
	if !hasSpecial {
		return false
	}

	return true
}

// GenerateCookieValue returns a 32-byte random, URL-safe string for session cookies.
// Panic on entropy failure to avoid issuing weak tokens.
func GenerateCookieValue() string {
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		panic("crypto/rand failed: " + err.Error())
	}

	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, 32)
	for i := range result {
		result[i] = letters[bytes[i]%byte(len(letters))]
	}

	return string(result)
}

// ======================
// Section: Sessions
// ======================

// IsValidSession returns true if the given session cookie exists and is not expired.
// Treats missing rows or DB errors as invalid (fail closed). Does not write to response.
func IsValidSession(cookieValue string) bool {
	db := database.CreateTable()
	defer db.Close()

	var expiration time.Time
	err := db.QueryRow(
		"SELECT expiration_date FROM Sessions WHERE cookie_value = ?",
		cookieValue,
	).Scan(&expiration)
	if err != nil {
		return false // not found, or some other DB error
	}
	return time.Now().Before(expiration)
}

// GetUserIDFromSession returns the user ID for a given session cookie
// Returns 0 if the session is invalid/expired; callers must handle 0 as unauthenticated.
func GetUserIDFromSession(cookieValue string) int {
	db := database.CreateTable()
	defer db.Close()

	var userID int
	err := db.QueryRow("SELECT user_id FROM Sessions WHERE cookie_value = ? AND expiration_date > datetime('now')", cookieValue).Scan(&userID)
	if err != nil {
		return 0
	}
	return userID
}

// GetUsernameFromSession returns the username for a given session cookie
// Returns empty string if session is invalid/expired or user lookup fails.
func GetUsernameFromSession(cookieValue string) string {
	db := database.CreateTable()
	defer db.Close()

	var username string
	err := db.QueryRow(`
		SELECT u.username 
		FROM Users u 
		JOIN Sessions s ON u.user_id = s.user_id 
		WHERE s.cookie_value = ? AND s.expiration_date > datetime('now')
	`, cookieValue).Scan(&username)
	if err != nil {
		return ""
	}
	return username
}

// CheckAuth is a middleware to check if user is authenticated
// Reads the "session" cookie, validates it, and returns (ok, userID, username).
// Early returns false if cookie is missing or session invalid; does not modify the response.
func CheckAuth(r *http.Request) (bool, int, string) {
	cookie, err := r.Cookie("session")
	if err != nil {
		return false, 0, ""
	}

	if !IsValidSession(cookie.Value) {
		return false, 0, ""
	}

	userID := GetUserIDFromSession(cookie.Value)
	username := GetUsernameFromSession(cookie.Value)

	return true, userID, username
}

// UpdateSessionUsername updates the username in the Sessions table
// Keeps session-linked displays in sync after username changes; no-op on DB error.
func UpdateSessionUsername(cookieValue string, newUsername string) {
	db := database.CreateTable()
	defer db.Close()

	// Ενημερώνουμε το username για το session ώστε να φαίνεται άμεσα στο frontend
	_, err := db.Exec(`
		UPDATE Users 
		SET username = ?
		WHERE user_id = (
			SELECT user_id 
			FROM Sessions 
			WHERE cookie_value = ?
		)
	`, newUsername, cookieValue)

	if err != nil {
		fmt.Println("Error updating session username:", err)
	}
}

// ======================
// Section: Formatting
// ======================

func FormatTimeAgo(t time.Time) string {
	duration := time.Since(t)

	switch {
	case duration < time.Minute:
		return "just now"
	case duration < time.Hour:
		minutes := int(duration.Minutes())
		if minutes == 1 {
			return "1 minute ago"
		}
		return fmt.Sprintf("%d minutes ago", int(duration.Minutes()))
	case duration < 24*time.Hour:
		hours := int(duration.Hours())
		if hours == 1 {
			return "1 hour ago"
		}
		return fmt.Sprintf("%d hours ago", int(duration.Hours()))
	case duration < 30*24*time.Hour:
		days := int(duration.Hours() / 24)
		if days == 1 {
			return "1 day ago"
		}
		return fmt.Sprintf("%d days ago", int(duration.Hours()/24))
	case duration < 12*30*24*time.Hour:
		months := int(duration.Hours() / (24 * 30))
		if months == 1 {
			return "1 month ago"
		}
		return fmt.Sprintf("%d months ago", int(duration.Hours()/(24*30)))
	default:
		years := int(duration.Hours() / (24 * 365))
		if years == 1 {
			return "1 year ago"
		}
		return fmt.Sprintf("%d years ago", int(duration.Hours()/(24*365)))
	}
}

// truncateText shortens text to a maximum length and adds "..." if truncated
func TruncateText(text string, maxLength int) string {
	if len(text) <= maxLength {
		return text
	}
	return text[:maxLength] + "..."
}

// IsAPIRequest checks if the path is an API endpoint
// Helps choose response format (JSON vs HTML) in routing wrappers.
func IsAPIRequest(path string) bool {
	return strings.HasPrefix(path, "/api/") ||
		strings.HasPrefix(path, "/auth/") ||
		strings.HasPrefix(path, "/ws") ||
		strings.HasPrefix(path, "/frontend/")
}
