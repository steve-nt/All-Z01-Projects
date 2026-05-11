package websocket

import (
	"log"
)

// NotificationData represents the data structure for a notification message
type NotificationData struct {
	NotificationID   int    `json:"notification_id,omitempty"`
	Type             string `json:"type"`              // follow_request, group_invitation, group_join_request, etc.
	Message          string `json:"message"`           // notification message text
	RelatedUserID    *int   `json:"related_user_id,omitempty"`
	RelatedGroupID   *int   `json:"related_group_id,omitempty"`
	RelatedPostID    *int   `json:"related_post_id,omitempty"`
	RelatedCommentID *int   `json:"related_comment_id,omitempty"`
	RelatedEventID   *int   `json:"related_event_id,omitempty"`
	IsRead           bool   `json:"is_read"`
	CreatedAt        string `json:"created_at,omitempty"`
}

// SendNotification sends a real-time notification to a user via WebSocket
// This should be called whenever a notification is created in the database
func SendNotification(userID int, notificationType string, message string, relatedUserID *int, relatedGroupID *int, relatedPostID *int, relatedCommentID *int, relatedEventID *int) {
	hub := GetGlobalHub()
	if hub == nil {
		log.Printf("Warning: WebSocket hub not initialized, cannot send notification to user %d", userID)
		return
	}

	notificationData := NotificationData{
		Type:             notificationType,
		Message:          message,
		RelatedUserID:    relatedUserID,
		RelatedGroupID:   relatedGroupID,
		RelatedPostID:    relatedPostID,
		RelatedCommentID: relatedCommentID,
		RelatedEventID:   relatedEventID,
		IsRead:           false,
	}

	wsMessage := &Message{
		Type:      "notification",
		Recipient: userID,
		Data:      notificationData,
	}

	hub.BroadcastToUser(userID, wsMessage)
	log.Printf("Sent real-time notification to user %d: type=%s", userID, notificationType)
}

// SendNotificationSimple is a convenience function for simple notifications
func SendNotificationSimple(userID int, notificationType string, message string) {
	SendNotification(userID, notificationType, message, nil, nil, nil, nil, nil)
}

// SendNotificationWithUser is a convenience function for notifications related to another user
func SendNotificationWithUser(userID int, notificationType string, message string, relatedUserID int) {
	SendNotification(userID, notificationType, message, &relatedUserID, nil, nil, nil, nil)
}

// SendNotificationWithGroup is a convenience function for notifications related to a group
func SendNotificationWithGroup(userID int, notificationType string, message string, relatedGroupID int) {
	SendNotification(userID, notificationType, message, nil, &relatedGroupID, nil, nil, nil)
}

// SendNotificationWithGroupAndUser is a convenience function for notifications with both group and user
func SendNotificationWithGroupAndUser(userID int, notificationType string, message string, relatedUserID int, relatedGroupID int) {
	SendNotification(userID, notificationType, message, &relatedUserID, &relatedGroupID, nil, nil, nil)
}

// SendNotificationWithEvent is a convenience function for event-related notifications
func SendNotificationWithEvent(userID int, notificationType string, message string, relatedUserID int, relatedGroupID int, relatedEventID int) {
	SendNotification(userID, notificationType, message, &relatedUserID, &relatedGroupID, nil, nil, &relatedEventID)
}

