package handlers

import (
	"bytes"
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"strings"

	"forum/db"
	"forum/sessions"

	"golang.org/x/crypto/bcrypt"
)

// Register shows the registration form (GET) and creates a new user (POST).
func Register(w http.ResponseWriter, r *http.Request) {
	switch r.Method {

	case http.MethodGet:
		// Render templates/register.html into the layout
		var buf bytes.Buffer
		if err := Templates.ExecuteTemplate(&buf, "register", nil); err != nil {
			RenderError(w, r, http.StatusInternalServerError, "Template render error.")
			return
		}
		RenderPage(w, r, "Register", template.HTML(buf.String()))
		return

	case http.MethodPost:
		_ = r.ParseForm()

		email := strings.TrimSpace(strings.ToLower(r.FormValue("email")))
		username := strings.TrimSpace(r.FormValue("username"))
		password := r.FormValue("password")

		// Safety net (browser should block this normally)
		if email == "" || username == "" || password == "" {
			RenderError(w, r, http.StatusBadRequest, "Please fill in all fields.")
			return
		}

		// Email already exists?
		if _, err := db.GetUserByEmail(email); err == nil {
			RenderError(w, r, http.StatusBadRequest, "That email is already registered. Try logging in instead.")
			return
		} else if err != sql.ErrNoRows {
			log.Println("GetUserByEmail error:", err)
			RenderError(w, r, http.StatusInternalServerError, "Something went wrong. Please try again.")
			return
		}

		// Username already exists?
		if _, err := db.GetUserByUsername(username); err == nil {
			RenderError(w, r, http.StatusBadRequest, "That username is already taken. Please choose another one.")
			return
		} else if err != sql.ErrNoRows {
			log.Println("GetUserByUsername error:", err)
			RenderError(w, r, http.StatusInternalServerError, "Something went wrong. Please try again.")
			return
		}

		// Hash password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			log.Println("bcrypt error:", err)
			RenderError(w, r, http.StatusInternalServerError, "Could not process password.")
			return
		}

		// Create user
		if err := db.CreateUser(email, username, string(hashedPassword)); err != nil {
			log.Println("CreateUser error:", err)
			RenderError(w, r, http.StatusBadRequest, "Email or username is already taken. Please try different values.")
			return
		}

		// Redirect to login page with a success flag.
		http.Redirect(w, r, "/login?registered=1", http.StatusSeeOther)
		return

	default:
		RenderError(w, r, http.StatusMethodNotAllowed, "Method not allowed.")
	}
}

// Login shows the login form (GET) and authenticates the user (POST).
func Login(w http.ResponseWriter, r *http.Request) {
	switch r.Method {

	case http.MethodGet:
		// Show a one-time success message after registration.
		pageData := map[string]interface{}{
			"RegisteredSuccess": r.URL.Query().Get("registered") == "1",
		}

		// Render templates/login.html into the layout
		var buf bytes.Buffer
		if err := Templates.ExecuteTemplate(&buf, "login", pageData); err != nil {
			RenderError(w, r, http.StatusInternalServerError, "Template render error.")
			return
		}
		RenderPage(w, r, "Login", template.HTML(buf.String()))
		return

	case http.MethodPost:
		_ = r.ParseForm()

		email := strings.TrimSpace(strings.ToLower(r.FormValue("email")))
		password := r.FormValue("password")

		if email == "" || password == "" {
			RenderError(w, r, http.StatusBadRequest, "Please enter email and password.")
			return
		}

		// Lookup user by email
		user, err := db.GetUserByEmail(email)
		if err != nil {
			RenderError(w, r, http.StatusUnauthorized, "Wrong email or password.")
			return
		}

		// Compare bcrypt hashes
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
			RenderError(w, r, http.StatusUnauthorized, "Wrong email or password.")
			return
		}

		// Create DB-backed session + cookie
		if err := sessions.CreateSession(w, r, user.ID); err != nil {
			log.Println("CreateSession error:", err)
			RenderError(w, r, http.StatusInternalServerError, "Could not create session.")
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
		return

	default:
		RenderError(w, r, http.StatusMethodNotAllowed, "Method not allowed.")
	}
}