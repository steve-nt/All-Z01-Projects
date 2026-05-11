package models

import (
	"errors"
	"forum/src/utils"
)

type Notifications []Notification

func GetNotificationsByUserId(userId int64) (Notifications, error) {
	var notifications Notifications
	rows, err := db.Query(`
	SELECT n.id, n.user_id, n.actor_id, n.type, n.post_id, comment_id, n.timestamp, n.read, u.username
	FROM notifications n
	JOIN users u ON u.id = n.actor_id
	WHERE user_id = ?
	ORDER BY n.timestamp DESC
	`, userId)
	if err != nil {
		err = errors.Join(utils.GetFunctionName(), err)
		return Notifications{}, err
	}
	defer rows.Close()
	for rows.Next() {
		var notification Notification
		err = rows.Scan(&notification.Id,
			&notification.UserId,
			&notification.ActorId,
			&notification.Type,
			&notification.PostId,
			&notification.CommentId,
			&notification.TimestampString,
			&notification.Read,
			&notification.Actor.Username)
		if err != nil {
			err = errors.Join(utils.GetFunctionName(), err)
			return Notifications{}, err
		}
		notifications = append(notifications, notification)
	}
	return notifications, nil
}
