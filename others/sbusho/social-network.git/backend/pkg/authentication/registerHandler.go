package authentication

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"social-network/backend/pkg/db/sqlite"
	"social-network/backend/utils"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.FileService("register.html", w, nil)
		return
	}

	// Parse multipart form to handle file uploads (32MB max memory)
	// This allows avatar upload during registration
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		// If multipart parsing fails, try regular form parsing (for forms without file upload)
		r.ParseForm()
	}

	db := sqlite.GetDB()

	email := strings.TrimSpace(r.FormValue("email"))
	pass := r.FormValue("password")
	confirm := r.FormValue("confirmPassword")
	first_name := strings.TrimSpace(r.FormValue("first_name"))
	last_name := strings.TrimSpace(r.FormValue("last_name"))
	date_of_birth := strings.TrimSpace(r.FormValue("date_of_birth"))
	nickname := strings.TrimSpace(r.FormValue("nickname"))
	about_me := strings.TrimSpace(r.FormValue("about_me"))
	is_public := strings.TrimSpace(r.FormValue("is_public"))
	is_active := strings.TrimSpace(r.FormValue("is_active"))

	// Handle avatar upload (optional)
	// Try to get avatar file - if it doesn't exist, that's fine (avatar is optional)
	var avatar_path string = "" // Will be set after user creation if avatar is uploaded
	var avatarFile multipart.File
	var avatarHeader *multipart.FileHeader
	avatarFile, avatarHeader, err = r.FormFile("avatar")
	if err != nil {
		// No avatar file provided - that's okay, it's optional
		avatarFile = nil
		avatarHeader = nil
	}

	formData := map[string]interface{}{
		"Email":       email,
		"FirstName":   first_name,
		"LastName":    last_name,
		"DateOfBirth": date_of_birth,
		"Nickname":    nickname,
		"AboutMe":     about_me,
		"IsPublic":    is_public,
		"IsActive":    is_active,
	}

	// Validate required fields
	if email == "" || first_name == "" || last_name == "" || date_of_birth == "" || pass == "" || confirm == "" {
		if wantsJSON(r) {
			writeJSON(w, http.StatusBadRequest, map[string]any{"error": "All required fields must be filled"})
			return
		}
		formData["Message"] = "All required fields must be filled"
		utils.FileService("register.html", w, formData)
		return
	}

	if pass != confirm {
		if wantsJSON(r) {
			writeJSON(w, http.StatusBadRequest, map[string]any{"error": "Passwords do not match"})
			return
		}
		formData["Message"] = "Passwords do not match"
		utils.FileService("register.html", w, formData)
		return
	}

	emailValid := utils.IsValidEmail(email)

	if !emailValid {
		if wantsJSON(r) {
			writeJSON(w, http.StatusBadRequest, map[string]any{"error": "Invalid email format"})
			return
		}
		formData["Message"] = "Invalid email format"
		utils.FileService("register.html", w, formData)
		return

	}

	passValid := utils.IsValidPassword(pass)

	if !passValid {
		if wantsJSON(r) {
			writeJSON(w, http.StatusBadRequest, map[string]any{"error": "Password must have: 8+ characters, 1 uppercase, 1 lowercase, 1 number, 1 symbol"})
			return
		}
		formData["Message"] = "Password must have: 8+ characters, 1 uppercase, 1 lowercase, 1 number, 1 symbol"
		utils.FileService("register.html", w, formData)
		return
	}

	// Check for duplicates
	var exists int
	err = db.QueryRow("SELECT COUNT(*) FROM Users WHERE email = ?", email).Scan(&exists)
	if err != nil {
		if wantsJSON(r) {
			writeJSON(w, http.StatusInternalServerError, map[string]any{"error": "Database error occurred"})
			return
		}
		formData["Message"] = "Database error occurred"
		utils.FileService("register.html", w, formData)
		return
	}

	if exists > 0 {
		if wantsJSON(r) {
			writeJSON(w, http.StatusConflict, map[string]any{"error": "This email is already registered"})
			return
		}
		formData["Message"] = "This email is already registered"
		utils.FileService("register.html", w, formData)
		return
	}

	// Convert is_public and is_active to boolean
	isPublicBool := is_public == "true" || is_public == "1" || is_public == "on"
	isActiveBool := is_active == "true" || is_active == "1" || is_active == "on"
	// Default to true if not provided
	if is_public == "" {
		isPublicBool = true
	}
	if is_active == "" {
		isActiveBool = true
	}

	// Hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		if wantsJSON(r) {
			writeJSON(w, http.StatusInternalServerError, map[string]any{"error": "Internal server error"})
			return
		}
		formData["Message"] = "Internal server error"
		utils.FileService("register.html", w, formData)
		return
	}

	// Insert user with all fields (required and optional)
	// Note: avatar_path will be set after user creation if avatar was uploaded
	_, err = db.Exec(`
		INSERT INTO Users (email, password_hash, first_name, last_name, date_of_birth, avatar_path, nickname, about_me, is_public, is_active) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		email, string(hash), first_name, last_name, date_of_birth,
		avatar_path, nickname, about_me, isPublicBool, isActiveBool)
	if err != nil {
		if wantsJSON(r) {
			writeJSON(w, http.StatusInternalServerError, map[string]any{"error": "Failed to create user account"})
			return
		}
		formData["Message"] = "Failed to create user account"
		utils.FileService("register.html", w, formData)
		return
	}

	// Get the newly created user ID
	var newUserID int
	err = db.QueryRow("SELECT user_id FROM Users WHERE email = ?", email).Scan(&newUserID)
	if err != nil {
		if wantsJSON(r) {
			writeJSON(w, http.StatusInternalServerError, map[string]any{"error": "Failed to retrieve user ID"})
			return
		}
		formData["Message"] = "Failed to retrieve user ID"
		utils.FileService("register.html", w, formData)
		return
	}

	// Handle avatar upload if provided (optional field)
	if avatarFile != nil && avatarHeader != nil {
		// Reset file pointer to beginning
		avatarFile.Seek(0, 0)

		// Process avatar upload using the helper function from avatar_image.go
		avatarPath, uploadErr := processAvatarUpload(avatarFile, avatarHeader, newUserID)
		if uploadErr == nil {
			// Update user with avatar path
			_, updateErr := db.Exec("UPDATE Users SET avatar_path = ? WHERE user_id = ?", avatarPath, newUserID)
			if updateErr != nil {
				// Log error but don't fail registration (avatar is optional)
				fmt.Printf("Warning: Could not update avatar path: %v\n", updateErr)
			}
		} else {
			// Log error but don't fail registration (avatar is optional)
			fmt.Printf("Warning: Could not upload avatar: %v\n", uploadErr)
		}
		// Close the file
		avatarFile.Close()
	}

	// Create welcome notification
	if newUserID > 0 {
		displayName := nickname
		if displayName == "" {
			displayName = first_name
		}
		message := fmt.Sprintf("Welcome to our community, %s! Start by creating your first post!", displayName)

		if _, err = sqlite.Insert(db, "Notifications", "(user_id, type, message, created_at)", newUserID, "welcome", message, time.Now()); err != nil {
			fmt.Printf("Warning: Could not create welcome notification: %v\n", err)
		}
	}

	if wantsJSON(r) {
		writeJSON(w, http.StatusCreated, map[string]any{"ok": true, "userID": newUserID})
		return
	}

	http.Redirect(w, r, "/login?success=registration", http.StatusSeeOther)
}
