package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"forum/db"
)

// CreateCategory handles category creation from the admin dashboard.
func CreateCategory(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		RenderError(w, r, http.StatusMethodNotAllowed, "Method not allowed.")
		return
	}

	_, ok := requireAdmin(w, r)
	if !ok {
		return
	}

	name := strings.TrimSpace(r.FormValue("name"))
	if name == "" {
		RenderError(w, r, http.StatusBadRequest, "Category name is required.")
		return
	}

	if err := db.CreateCategory(name); err != nil {
		RenderError(w, r, http.StatusInternalServerError, "Could not create category.")
		return
	}

	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

// DeleteCategory handles category deletion from the admin dashboard.
func DeleteCategory(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		RenderError(w, r, http.StatusMethodNotAllowed, "Method not allowed.")
		return
	}

	_, ok := requireAdmin(w, r)
	if !ok {
		return
	}

	categoryID, err := strconv.Atoi(r.FormValue("category_id"))
	if err != nil || categoryID <= 0 {
		RenderError(w, r, http.StatusBadRequest, "Invalid category.")
		return
	}

	if err := db.DeleteCategory(categoryID); err != nil {
		RenderError(w, r, http.StatusInternalServerError, "Could not delete category.")
		return
	}

	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}