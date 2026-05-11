package handlers

import (
	"bytes"
	"html/template"
	"net/http"

	"forum/db"
)

// RequestModerator shows the request form and lets a normal user ask for moderator access.
func RequestModerator(w http.ResponseWriter, r *http.Request) {
	userID, ok := requireLogin(w, r)
	if !ok {
		return
	}

	role, err := db.GetUserRole(userID)
	if err != nil {
		RenderError(w, r, http.StatusInternalServerError, "Could not verify user role.")
		return
	}

	// Only normal users should request promotion.
	if role != "user" {
		RenderError(w, r, http.StatusForbidden, "Only normal users can request moderator access.")
		return
	}

	switch r.Method {
	case http.MethodGet:
		hasPending, err := db.HasPendingModeratorRequest(userID)
		if err != nil {
			RenderError(w, r, http.StatusInternalServerError, "Could not load request status.")
			return
		}

		pageData := map[string]interface{}{
			"HasPending": hasPending,
		}

		var buf bytes.Buffer
		if err := Templates.ExecuteTemplate(&buf, "request_moderator", pageData); err != nil {
			RenderError(w, r, http.StatusInternalServerError, "Template render error.")
			return
		}

		RenderPage(w, r, "Request Moderator Access", template.HTML(buf.String()))
		return

	case http.MethodPost:
		hasPending, err := db.HasPendingModeratorRequest(userID)
		if err != nil {
			RenderError(w, r, http.StatusInternalServerError, "Could not validate request status.")
			return
		}
		if hasPending {
			RenderError(w, r, http.StatusBadRequest, "You already have a pending moderator request.")
			return
		}

		if err := db.CreateModeratorRequest(userID); err != nil {
			RenderError(w, r, http.StatusInternalServerError, "Could not create moderator request.")
			return
		}

		// AUDIT: notify all admins that a moderator request was submitted.
		adminIDs, err := db.GetAdminUserIDs()
		if err == nil {
			for _, adminID := range adminIDs {
				_ = db.CreateCustomNotification(
					adminID,
					userID,
					0,
					"moderator_request_submitted",
					"A user submitted a moderator access request.",
				)
			}
		}

		http.Redirect(w, r, "/request-moderator", http.StatusSeeOther)
		return

	default:
		RenderError(w, r, http.StatusMethodNotAllowed, "Method not allowed.")
		return
	}
}