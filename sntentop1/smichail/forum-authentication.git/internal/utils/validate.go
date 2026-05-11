package utils

import (
	"encoding/json"
	"html"
	"log"
	"net/http"
	"regexp"
	"strings"
)

// SanitizedUserInput holds cleaned values for a new user
type SanitizedUserInput struct {
	Mail           string
	Username       string
	Password       string
	RepeatPassword string
	Role           string
}

// ValidationErrors holds all validation errors
type ValidationErrors struct {
	Errors []string `json:"errors"`
}

var (
	emailRegex    = regexp.MustCompile(`^[A-Za-z0-9._%+\-]+@[A-Za-z0-9.\-]+\.[A-Za-z]{2,}$`)
	usernameRegex = regexp.MustCompile(`^[A-Za-z0-9_]{3,20}$`)
	hasUpper      = regexp.MustCompile(`[A-Z]`)
	hasLower      = regexp.MustCompile(`[a-z]`)
	hasDigit      = regexp.MustCompile(`[0-9]`)
	hasSpecial    = regexp.MustCompile(`[^A-Za-z0-9]`)
)

// SanitizeAndValidateNewUser trims whitespace, escapes HTML, and validates fields.
// Returns sanitized values and a slice of all validation errors encountered.
func SanitizeAndValidateNewUser(mail, username, password, repeat_password, role string) (SanitizedUserInput, []string) {
	var errors []string

	clean := SanitizedUserInput{
		Mail:           strings.TrimSpace(mail),
		Username:       strings.TrimSpace(username),
		Password:       strings.TrimSpace(password),
		RepeatPassword: strings.TrimSpace(repeat_password),
		Role:           strings.TrimSpace(role),
	}

	// Email format (only check if not empty)
	if clean.Mail != "" && !emailRegex.MatchString(clean.Mail) {
		errors = append(errors, "Invalid email format")
	}

	// Username rules: 3-20 chars, alnum and underscore only (only check if not empty)
	if clean.Username != "" && !usernameRegex.MatchString(clean.Username) {
		errors = append(errors, "Username must be 3-20 characters and contain only letters, numbers, and underscores")
	}

	// Password rules (only check if not empty)
	if clean.Password != "" {
		if len(clean.Password) < 8 {
			errors = append(errors, "Password must be at least 8 characters")
		}
		if !hasUpper.MatchString(clean.Password) {
			errors = append(errors, "Password must contain at least one uppercase letter")
		}
		if !hasLower.MatchString(clean.Password) {
			errors = append(errors, "Password must contain at least one lowercase letter")
		}
		if !hasDigit.MatchString(clean.Password) {
			errors = append(errors, "Password must contain at least one number")
		}
		if !hasSpecial.MatchString(clean.Password) {
			errors = append(errors, "Password must contain at least one special character")
		}
	}

	if clean.Password != clean.RepeatPassword {
		errors = append(errors, "Passwords do not match")
	}

	// Limit role to known values, default to "user"
	switch strings.ToLower(clean.Role) {
	case "", "user":
		clean.Role = "user"
	case "admin":
		clean.Role = "admin"
	default:
		clean.Role = "user"
	}

	// Escape any HTML to prevent injection when echoed
	clean.Mail = html.EscapeString(clean.Mail)
	clean.Username = html.EscapeString(clean.Username)
	// Do not escape password characters; pass through to hasher as-is after trim

	return clean, errors
}

// ErrorResponse sends a standardized JSON error response
func ErrorResponseSignup(w http.ResponseWriter, errors []string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := map[string]interface{}{
		"success": false,
		"status":  statusCode,
		"errors":  errors,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Println("ERROR encoding JSON response:", err)
	}
}

// SuccessResponse sends a standardized JSON success response
func SuccessResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := map[string]interface{}{
		"success": true,
		"status":  statusCode,
	}

	// If data is a string, treat it as a message
	if message, ok := data.(string); ok {
		response["message"] = message
	} else {
		response["data"] = data
	}

	json.NewEncoder(w).Encode(response)
}
