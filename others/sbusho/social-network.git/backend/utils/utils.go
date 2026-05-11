// Package utils provides shared helpers for templates, auth, validation, and formatting.
package utils

import (
	"crypto/rand"
	"fmt"
	"html/template"
	"net/http"
	"regexp"
	"social-network/backend/pkg/db/sqlite"
	"strings"
	"time"
)

// ===== Templates =====

// TemplateData holds data to pass to templates.
type TemplateData struct {
	IsLoggedIn bool
	Nickname   string
	UserID     int
	Message    string
	Error      string
	Data       interface{}
}

// FileService renders a template without authentication context.
func FileService(filename string, w http.ResponseWriter, data any) {
	tmpl, err := template.ParseFiles("frontend/templates/" + filename)
	if err != nil {
		panic("Template error: " + err.Error())
	}
	tmpl.Execute(w, data)
}

// FileServiceWithAuth serves templates with authentication context.
func FileServiceWithAuth(filename string, w http.ResponseWriter, r *http.Request, data interface{}) {
	templateData := &TemplateData{
		Data: data,
	}

	// Check if user is logged in
	if cookie, err := r.Cookie("session"); err == nil && IsValidSession(cookie.Value) {
		templateData.IsLoggedIn = true
		templateData.UserID = GetUserIDFromSession(cookie.Value)
		templateData.Nickname = GetNicknameFromSession(cookie.Value)
	}

	tmpl, err := template.ParseFiles("frontend/templates/" + filename)
	if err != nil {
		http.Error(w, "Template error: "+err.Error(), 500)
		return
	}
	tmpl.Execute(w, templateData)
}

// ===== Validation =====

// IsValidEmail performs a basic email format check.
func IsValidEmail(email string) bool {
	// Trim whitespace
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

// TestBasicEmailValidation prints a simple email validation demo.
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

// IsValidPassword enforces a basic complexity policy.
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

// ===== Auth Helpers =====

// GenerateCookieValue creates a random session cookie value.
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

// IsValidSession returns true if the given session cookie exists and is not expired.
func IsValidSession(cookieValue string) bool {
	db := sqlite.GetDB()

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

// GetUserIDFromSession returns the user ID for a given session cookie.
func GetUserIDFromSession(cookieValue string) int {
	db := sqlite.GetDB()

	var userID int
	err := db.QueryRow("SELECT user_id FROM Sessions WHERE cookie_value = ? AND expiration_date > datetime('now')", cookieValue).Scan(&userID)
	if err != nil {
		return 0
	}
	return userID
}

// GetNicknameFromSession returns a display name for the user.
// Since nickname is optional, it falls back to email if nickname is empty.
func GetNicknameFromSession(cookieValue string) string {
	db := sqlite.GetDB()

	var displayName string
	// COALESCE returns the first non-NULL value: nickname if set, otherwise first_name, otherwise email
	err := db.QueryRow(`
		SELECT COALESCE(NULLIF(u.nickname, ''), u.email) 
		FROM Users u 
		JOIN Sessions s ON u.user_id = s.user_id 
		WHERE s.cookie_value = ? AND s.expiration_date > datetime('now')
	`, cookieValue).Scan(&displayName)
	if err != nil {
		return ""
	}
	return displayName
}

// CheckAuth checks whether the request has a valid session cookie.
func CheckAuth(r *http.Request) (bool, int, string) {
	cookie, err := r.Cookie("session")
	if err != nil {
		return false, 0, ""
	}

	if !IsValidSession(cookie.Value) {
		return false, 0, ""
	}

	userID := GetUserIDFromSession(cookie.Value)
	nickname := GetNicknameFromSession(cookie.Value)

	return true, userID, nickname
}

// UpdateSessionNickname updates the nickname in the Users table for the given session.
// Note: nickname is optional, so this allows users to set/update their nickname.
func UpdateSessionNickname(cookieValue string, newNickname string) {
	db := sqlite.GetDB()

	_, err := db.Exec(`
		UPDATE Users 
		SET nickname = ?
		WHERE user_id = (
			SELECT user_id 
			FROM Sessions 
			WHERE cookie_value = ?
		)
	`, newNickname, cookieValue)

	if err != nil {
		fmt.Println("Error updating user nickname:", err)
	}
}

// ===== Formatting =====

// FormatTimeAgo renders a human-readable "time ago" string.
func FormatTimeAgo(t time.Time) string {
	duration := time.Since(t)

	switch {
	case duration < time.Minute:
		return "just now"
	case duration < time.Hour:
		return fmt.Sprintf("%d minutes ago", int(duration.Minutes()))
	case duration < 24*time.Hour:
		return fmt.Sprintf("%d hours ago", int(duration.Hours()))
	case duration < 30*24*time.Hour:
		return fmt.Sprintf("%d days ago", int(duration.Hours()/24))
	case duration < 12*30*24*time.Hour:
		return fmt.Sprintf("%d months ago", int(duration.Hours()/(24*30)))
	default:
		return fmt.Sprintf("%d years ago", int(duration.Hours()/(24*365)))
	}
}

// TruncateText shortens text to a maximum length and adds "..." if truncated.
func TruncateText(text string, maxLength int) string {
	if len(text) <= maxLength {
		return text
	}
	return text[:maxLength] + "..."
}
