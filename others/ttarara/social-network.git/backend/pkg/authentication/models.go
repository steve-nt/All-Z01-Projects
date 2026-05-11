package authentication

import "time"

// User represents a user in the system
// NOTE: This struct is optional for Part 1 - current implementation works without it
// Creating it for better code organization and type safety
// Can be used in Part 2 (Profiles) and future parts
type User struct {
	UserID      int       `json:"user_id" db:"user_id"`
	Email       string    `json:"email" db:"email"`
	PasswordHash string   `json:"-" db:"password_hash"` // Never expose password hash in JSON
	FirstName   string    `json:"first_name" db:"first_name"`
	LastName    string    `json:"last_name" db:"last_name"`
	DateOfBirth time.Time `json:"date_of_birth" db:"date_of_birth"`
	AvatarPath  *string   `json:"avatar_path,omitempty" db:"avatar_path"` // Pointer for nullable
	Nickname    *string   `json:"nickname,omitempty" db:"nickname"`        // Pointer for nullable
	AboutMe     *string   `json:"about_me,omitempty" db:"about_me"`        // Pointer for nullable
	IsPublic    bool      `json:"is_public" db:"is_public"`
	IsActive    bool      `json:"is_active" db:"is_active"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

// Session represents a user session
// NOTE: This struct is optional for Part 1 - current implementation works without it
type Session struct {
	SessionID    int       `json:"session_id" db:"session_id"`
	UserID       int       `json:"user_id" db:"user_id"`
	CookieValue  string    `json:"-" db:"cookie_value"` // Never expose cookie value in JSON
	ExpirationDate time.Time `json:"expiration_date" db:"expiration_date"`
}

// RegisterRequest represents registration form data
// NOTE: Optional - can make validation and handling cleaner
type RegisterRequest struct {
	Email       string `json:"email"`
	Password    string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	DateOfBirth string `json:"date_of_birth"` // String because it comes from form
	Nickname    string `json:"nickname"`
	AboutMe     string `json:"about_me"`
	IsPublic    string `json:"is_public"` // String because it comes from form
}

// LoginRequest represents login form data
// NOTE: Optional - can make validation cleaner
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

