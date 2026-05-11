package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"forum/db"
)

// UpdateUserRole changes a user's role from the admin dashboard.
func UpdateUserRole(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		RenderError(w, r, http.StatusMethodNotAllowed, "Method not allowed.")
		return
	}

	adminUserID, ok := requireAdmin(w, r)
	if !ok {
		return
	}

	targetUserID, err := strconv.Atoi(r.FormValue("user_id"))
	if err != nil || targetUserID <= 0 {
		RenderError(w, r, http.StatusBadRequest, "Invalid user.")
		return
	}

	role := strings.TrimSpace(r.FormValue("role"))
	if role == "" {
		RenderError(w, r, http.StatusBadRequest, "Role is required.")
		return
	}

	// AUDIT: avoid accidentally demoting the last active admin through the UI flow.
	if targetUserID == adminUserID && role != "admin" {
		RenderError(w, r, http.StatusBadRequest, "You cannot remove your own admin role.")
		return
	}

	if err := db.UpdateUserRole(targetUserID, role); err != nil {
		RenderError(w, r, http.StatusInternalServerError, "Could not update user role.")
		return
	}

	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

// ApproveModeratorRequest handles admin approval of a moderator request.
func ApproveModeratorRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		RenderError(w, r, http.StatusMethodNotAllowed, "Method not allowed.")
		return
	}

	adminUserID, ok := requireAdmin(w, r)
	if !ok {
		return
	}

	requestID, err := strconv.Atoi(r.FormValue("request_id"))
	if err != nil || requestID <= 0 {
		RenderError(w, r, http.StatusBadRequest, "Invalid request.")
		return
	}

	if err := db.ApproveModeratorRequest(requestID, adminUserID); err != nil {
		RenderError(w, r, http.StatusInternalServerError, "Could not approve moderator request.")
		return
	}

	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

// RejectModeratorRequest handles admin rejection of a moderator request.
func RejectModeratorRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		RenderError(w, r, http.StatusMethodNotAllowed, "Method not allowed.")
		return
	}

	adminUserID, ok := requireAdmin(w, r)
	if !ok {
		return
	}

	requestID, err := strconv.Atoi(r.FormValue("request_id"))
	if err != nil || requestID <= 0 {
		RenderError(w, r, http.StatusBadRequest, "Invalid request.")
		return
	}

	if err := db.RejectModeratorRequest(requestID, adminUserID); err != nil {
		RenderError(w, r, http.StatusInternalServerError, "Could not reject moderator request.")
		return
	}

	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}