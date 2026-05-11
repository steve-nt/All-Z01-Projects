package db

// Notification represents one item in the notifications list.
type Notification struct {
	ID            int
	UserID        int
	ActorUserID   int
	ActorUsername string
	PostID        int
	PostTitle     string
	Type          string
	Message       string
	IsRead        bool
	CreatedAt     string
}

// CreateNotification inserts a new notification.
// Self-notifications are ignored to avoid noisy UX.
func CreateNotification(userID, actorUserID, postID int, notifType string) error {
	if userID <= 0 || actorUserID <= 0 || notifType == "" {
		return nil
	}

	// Do not notify users about their own actions.
	if userID == actorUserID {
		return nil
	}

	_, err := DB.Exec(
		`INSERT INTO notifications (user_id, actor_user_id, post_id, type, message)
		 VALUES (?, ?, ?, ?, '')`,
		userID, actorUserID, nullablePostID(postID), notifType,
	)
	return err
}

// CreateCustomNotification inserts a notification with a custom text message.
// AUDIT: moderation workflows need notifications that are not always tied to the
// standard post liked/commented/disliked event set.
func CreateCustomNotification(userID, actorUserID, postID int, notifType, message string) error {
	if userID <= 0 || actorUserID <= 0 || notifType == "" {
		return nil
	}

	// Do not notify users about their own actions.
	if userID == actorUserID {
		return nil
	}

	_, err := DB.Exec(
		`INSERT INTO notifications (user_id, actor_user_id, post_id, type, message)
		 VALUES (?, ?, ?, ?, ?)`,
		userID, actorUserID, nullablePostID(postID), notifType, message,
	)
	return err
}

// CountUnreadNotifications returns the number of unread notifications for a user.
func CountUnreadNotifications(userID int) (int, error) {
	var n int
	err := DB.QueryRow(
		`SELECT COUNT(*) FROM notifications WHERE user_id = ? AND is_read = 0`,
		userID,
	).Scan(&n)
	return n, err
}

// GetNotificationsByUser returns notifications newest first.
// It also joins actor username and post title for display.
func GetNotificationsByUser(userID int) ([]Notification, error) {
	rows, err := DB.Query(`
		SELECT
			n.id,
			n.user_id,
			n.actor_user_id,
			u.username,
			IFNULL(n.post_id, 0),
			IFNULL(p.title, ''),
			n.type,
			IFNULL(n.message, ''),
			n.is_read,
			n.created_at
		FROM notifications n
		JOIN users u ON u.id = n.actor_user_id
		LEFT JOIN posts p ON p.id = n.post_id
		WHERE n.user_id = ?
		ORDER BY n.created_at DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []Notification

	for rows.Next() {
		var n Notification
		var isReadInt int

		if err := rows.Scan(
			&n.ID,
			&n.UserID,
			&n.ActorUserID,
			&n.ActorUsername,
			&n.PostID,
			&n.PostTitle,
			&n.Type,
			&n.Message,
			&isReadInt,
			&n.CreatedAt,
		); err != nil {
			return nil, err
		}

		n.IsRead = isReadInt == 1
		out = append(out, n)
	}

	return out, rows.Err()
}

// MarkAllNotificationsRead marks all notifications as read for a given user.
func MarkAllNotificationsRead(userID int) error {
	_, err := DB.Exec(`UPDATE notifications SET is_read = 1 WHERE user_id = ?`, userID)
	return err
}

// nullablePostID keeps backward-compatible notification inserts while allowing
// generic moderation notifications without a related post link.
func nullablePostID(postID int) interface{} {
	if postID <= 0 {
		return nil
	}
	return postID
}