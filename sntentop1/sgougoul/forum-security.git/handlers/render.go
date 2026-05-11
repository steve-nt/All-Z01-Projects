package handlers

import (
	"bytes"
	"html/template"
	"log"
	"net/http"

	"forum/db"
	"forum/sessions"
)

// RenderPage renders the global layout and injects page content.
// It also injects login status so the navbar can switch links.
func RenderPage(w http.ResponseWriter, r *http.Request, title string, content template.HTML) {
	userID, loggedIn := sessions.GetUserID(r)

	unread := 0
	username := ""
	role := ""

	if loggedIn {
		if n, err := db.CountUnreadNotifications(userID); err == nil {
			unread = n
		}

		// Load current username for the persistent welcome message.
		if u, err := db.GetUserByID(userID); err == nil {
			username = u.Username
			role = u.Role
		}
	}

	data := map[string]interface{}{
		"Title":               title,
		"LoggedIn":            loggedIn,
		"Username":            username,
		"Role":                role,
		"UnreadNotifications": unread,
		"Content":             content,
	}

	if err := Templates.ExecuteTemplate(w, "layout", data); err != nil {
		log.Println("RenderPage template error:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// RenderError sends a styled error page with the correct HTTP status.
// It renders templates/error.html and wraps it with the global layout.
func RenderError(w http.ResponseWriter, r *http.Request, status int, message string) {
	title := http.StatusText(status)
	if title == "" {
		title = "Error"
	}

	_, loggedIn := sessions.GetUserID(r)

	// Write status before writing body
	w.WriteHeader(status)

	pageData := map[string]interface{}{
		"Title":    title,
		"Message":  message,
		"LoggedIn": loggedIn,
	}

	var buf bytes.Buffer
	if err := Templates.ExecuteTemplate(&buf, "error", pageData); err != nil {
		http.Error(w, title+": "+message, status)
		return
	}

	RenderPage(w, r, title, template.HTML(buf.String()))
}