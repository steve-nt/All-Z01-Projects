package handlers

import (
	"fmt"
	"net/http"
	"realtimeforum/internals/database"
	"realtimeforum/internals/utils"
	"strconv"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	// GET → show the form
	if r.Method != http.MethodPost {
		utils.FileService("register.html", w, nil)
		return
	}

	db := database.CreateTable()
	defer db.Close()

	//------------------------------------------
	// 1. Collect form data
	//------------------------------------------

	username := strings.TrimSpace(r.FormValue("username"))
	firstName := strings.TrimSpace(r.FormValue("first_name"))
	lastName := strings.TrimSpace(r.FormValue("last_name"))
	ageStr := strings.TrimSpace(r.FormValue("age"))
	gender := strings.TrimSpace(r.FormValue("gender"))
	email := strings.TrimSpace(r.FormValue("email"))
	pass := r.FormValue("password")
	confirm := r.FormValue("confirm_password")

	formData := map[string]interface{}{
		"Username":  username,
		"Email":     email,
		"FirstName": firstName,
		"LastName":  lastName,
		"Age":       ageStr,
		"Gender":    gender,
	}

	//------------------------------------------
	// 2. Basic validation
	//------------------------------------------

	if username == "" || firstName == "" || lastName == "" || ageStr == "" ||
		gender == "" || email == "" || pass == "" || confirm == "" {

		formData["Message"] = "All fields are required."
		utils.FileService("register.html", w, formData)
		return
	}

	// Convert age
	age, err := strconv.Atoi(ageStr)
	if err != nil || age < 13 {
		formData["Message"] = "You must be at least 13 years old."
		utils.FileService("register.html", w, formData)
		return
	}

	if pass != confirm {
		formData["Message"] = "Passwords do not match."
		utils.FileService("register.html", w, formData)
		return
	}

	if !utils.IsValidEmail(email) {
		formData["Message"] = "Invalid email format."
		utils.FileService("register.html", w, formData)
		return
	}

	if !utils.IsValidPassword(pass) {
		formData["Message"] = "Password must have 8+ characters, uppercase, lowercase, number, and symbol."
		utils.FileService("register.html", w, formData)
		return
	}

	//------------------------------------------
	// 3. Check duplicates (email / username)
	//------------------------------------------

	var exists int

	err = db.QueryRow("SELECT COUNT(*) FROM Users WHERE email = ?", email).Scan(&exists)
	if err != nil {
		formData["Message"] = "Database error occurred."
		utils.FileService("register.html", w, formData)
		return
	}
	if exists > 0 {
		formData["Message"] = "This email is already registered."
		utils.FileService("register.html", w, formData)
		return
	}

	err = db.QueryRow("SELECT COUNT(*) FROM Users WHERE username = ?", username).Scan(&exists)
	if err != nil {
		formData["Message"] = "Database error occurred."
		utils.FileService("register.html", w, formData)
		return
	}
	if exists > 0 {
		formData["Message"] = "This username is already taken."
		utils.FileService("register.html", w, formData)
		return
	}

	//------------------------------------------
	// 4. Hash password
	//------------------------------------------

	hash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		formData["Message"] = "Internal server error."
		utils.FileService("register.html", w, formData)
		return
	}

	//------------------------------------------
	// 5. Insert new user
	//------------------------------------------

	_, err = db.Exec(`
		INSERT INTO Users 
			(username, age, gender, first_name, last_name, email, password_hash)
		VALUES 
			(?, ?, ?, ?, ?, ?, ?)
	`, username, age, gender, firstName, lastName, email, string(hash))

	if err != nil {
		formData["Message"] = "Failed to create user account."
		utils.FileService("register.html", w, formData)
		return
	}

	//------------------------------------------
	// 6. Load new user ID
	//------------------------------------------

	var newUserID int
	db.QueryRow("SELECT user_id FROM Users WHERE email = ?", email).Scan(&newUserID)

	//------------------------------------------
	// 7. Create welcome notification
	//------------------------------------------

	if newUserID > 0 {
		title := "Welcome to Tech Talk! 🤖 💻"
		message := fmt.Sprintf(
			"Welcome %s! Glad to have you in our tech community. Start exploring, posting, and connecting!",
			username,
		)
		CreateNotification(newUserID, "system", title, message, nil, nil, nil)
	}

	//------------------------------------------
	// 8. Redirect to login with success flag
	//------------------------------------------

	http.Redirect(w, r, "/login?success=registration", http.StatusSeeOther)
}
