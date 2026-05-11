package handlers

import (
	"bytes"
	"html/template"
	"net/http"
	"strings"

	"forum/db"
)

// AdminDashboard renders the administrator control page.
// Only admins may access this page.
func AdminDashboard(w http.ResponseWriter, r *http.Request) {
	_, ok := requireAdmin(w, r)
	if !ok {
		return
	}

	users, err := db.GetAllUsers()
	if err != nil {
		RenderError(w, r, http.StatusInternalServerError, "Could not load users.")
		return
	}

	categories, err := db.GetAllCategories()
	if err != nil {
		RenderError(w, r, http.StatusInternalServerError, "Could not load categories.")
		return
	}

	requests, err := db.GetAllModeratorRequests()
	if err != nil {
		RenderError(w, r, http.StatusInternalServerError, "Could not load moderator requests.")
		return
	}

	reports, err := db.GetAllReports()
	if err != nil {
		RenderError(w, r, http.StatusInternalServerError, "Could not load reports.")
		return
	}

	for i := range requests {
		requests[i].CreatedAt = FormatDisplayTime(requests[i].CreatedAt)
		if strings.TrimSpace(requests[i].ReviewedAt) != "" {
			requests[i].ReviewedAt = FormatDisplayTime(requests[i].ReviewedAt)
		}
	}

	for i := range reports {
		reports[i].CreatedAt = FormatDisplayTime(reports[i].CreatedAt)
		if strings.TrimSpace(reports[i].ReviewedAt) != "" {
			reports[i].ReviewedAt = FormatDisplayTime(reports[i].ReviewedAt)
		}
	}

	pageData := map[string]interface{}{
		"Users":      users,
		"Categories": categories,
		"Requests":   requests,
		"Reports":    reports,
	}

	var buf bytes.Buffer
	if err := Templates.ExecuteTemplate(&buf, "admin", pageData); err != nil {
		RenderError(w, r, http.StatusInternalServerError, "Template render error.")
		return
	}

	RenderPage(w, r, "Admin Dashboard", template.HTML(buf.String()))
}