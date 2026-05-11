package handlers

import (
	"net/http"
	"realtimeforum/internals/utils"
)

// HomeHandler handles the main homepage route - serves SPA
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	// For SPA, serve index.html for all GET requests that aren't API or static files
	if r.Method == http.MethodGet && !utils.IsAPIRequest(r.URL.Path) {
		http.ServeFile(w, r, "frontend/index.html")
		return
	}
	
	// If it's not a GET request or is an API request, let other handlers deal with it
	NotFoundHandler(w, r)
}

// ViewPostHandler serves the view-post page
func ViewPostHandler(w http.ResponseWriter, r *http.Request) {
	utils.FileService("view-post.html", w, nil)
}

// CategoriesPageHandler serves the categories page
func CategoriesPageHandler(w http.ResponseWriter, r *http.Request) {
	utils.FileService("categories.html", w, nil)
}

// Static page handlers
func AboutHandler(w http.ResponseWriter, r *http.Request) {
	utils.FileService("about.html", w, nil)
}

func TermsHandler(w http.ResponseWriter, r *http.Request) {
	utils.FileService("Terms&Conditions.html", w, nil)
}

func ForgotPasswordPageHandler(w http.ResponseWriter, r *http.Request) {
	utils.FileService("forgot-password.html", w, nil)
}

func ProfilePageHandler(w http.ResponseWriter, r *http.Request) {
	if cookie, err := r.Cookie("session"); err == nil && utils.IsValidSession(cookie.Value) {
		utils.FileService("profile.html", w, nil)
	} else {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}

func NotificationsPageHandler(w http.ResponseWriter, r *http.Request) {
	if cookie, err := r.Cookie("session"); err == nil && utils.IsValidSession(cookie.Value) {
		utils.FileService("notifications.html", w, nil)
	} else {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}
