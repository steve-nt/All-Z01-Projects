package notifications

import (
	"net/http"
	"social-network/backend/middleware"
)

// SetupNotificationRoutes registers notification API endpoints
func SetupNotificationRoutes() {
	// GET /api/notifications - Get all notifications for the current user
	// Query params: ?unread_only=true&limit=50
	http.HandleFunc("/api/notifications", middleware.WrapHandler(middleware.RequireAuthJSON(NotificationsHandler)))

	// POST /api/notifications/read - Mark notification(s) as read
	// Body: {"notification_id": 123} or {"mark_all_read": true}
	http.HandleFunc("/api/notifications/read", middleware.WrapHandler(middleware.RequireAuthJSON(MarkReadHandler)))
}

