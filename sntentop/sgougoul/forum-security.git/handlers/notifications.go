package handlers

import (
	"bytes"
	"html/template"
	"net/http"

	"forum/db"
	"forum/sessions"
)

// NotificationView prepares DB notifications for display.
type NotificationView struct {
	ActorUsername string
	PostID        int
	PostTitle     string
	Message       string
	CreatedAt     string
	HasPostLink   bool
}

// Notifications renders the user's notifications page.
// Viewing this page also marks all notifications as read.
func Notifications(w http.ResponseWriter, r *http.Request) {
	userID, ok := sessions.GetUserID(r)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	items, err := db.GetNotificationsByUser(userID)
	if err != nil {
		RenderError(w, r, http.StatusInternalServerError, "Could not load notifications.")
		return
	}

	views := make([]NotificationView, 0, len(items))
	for _, n := range items {
		msg := n.Message

		// AUDIT: keep existing standard notifications working, while also
		// supporting custom moderation-related messages from the DB.
		if msg == "" {
			msg = "interacted with your post"
			switch n.Type {
			case "post_liked":
				msg = "liked your post"
			case "post_disliked":
				msg = "disliked your post"
			case "post_commented":
				msg = "commented on your post"
			case "post_pending_review":
				msg = "submitted a post that is waiting for moderation review"
			case "moderator_request_submitted":
				msg = "submitted a moderator access request"
			case "report_answered":
				msg = "received an admin response to a report"
			}
		}

		views = append(views, NotificationView{
			ActorUsername: n.ActorUsername,
			PostID:        n.PostID,
			PostTitle:     n.PostTitle,
			Message:       msg,
			CreatedAt:     FormatDisplayTime(n.CreatedAt),
			HasPostLink:   n.PostID > 0 && n.PostTitle != "",
		})
	}

	// Mark as read after loading.
	_ = db.MarkAllNotificationsRead(userID)

	pageData := map[string]interface{}{
		"Notifications": views,
	}

	var buf bytes.Buffer
	if err := Templates.ExecuteTemplate(&buf, "notifications", pageData); err != nil {
		RenderError(w, r, http.StatusInternalServerError, "Template render error.")
		return
	}

	RenderPage(w, r, "Notifications", template.HTML(buf.String()))
}