package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"forum/db"
)

// ResolveReport handles the admin response to a moderator report.
func ResolveReport(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		RenderError(w, r, http.StatusMethodNotAllowed, "Method not allowed.")
		return
	}

	adminUserID, ok := requireAdmin(w, r)
	if !ok {
		return
	}

	reportID, err := strconv.Atoi(r.FormValue("report_id"))
	if err != nil || reportID <= 0 {
		RenderError(w, r, http.StatusBadRequest, "Invalid report.")
		return
	}

	response := strings.TrimSpace(r.FormValue("response"))
	if response == "" {
		RenderError(w, r, http.StatusBadRequest, "Admin response is required.")
		return
	}

	// AUDIT: resolve the report first.
	if err := db.ResolveReport(reportID, adminUserID, response); err != nil {
		RenderError(w, r, http.StatusInternalServerError, "Could not resolve report.")
		return
	}

	// AUDIT: notify the moderator that the admin responded to the report.
	reports, err := db.GetAllReports()
	if err == nil {
		for _, rep := range reports {
			if rep.ID == reportID {
				_ = db.CreateCustomNotification(
					rep.ReporterID,
					adminUserID,
					rep.PostID,
					"report_answered",
					"An administrator responded to your report: "+response,
				)
				break
			}
		}
	}

	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}