package models

import "forum/src/utils"

type Notification struct {
	Id        int64
	UserId    int64
	ActorId   int64
	Actor     User
	Type      string
	Post      Post
	PostId    int64
	Comment	  Comment
	CommentId int64
	TimestampString string
	Read      bool
}

func (n *Notification) Add() error {
	query := `INSERT INTO notifications (user_id, actor_id, type, post_id, comment_id, timestamp) VALUES (?, ?, ?, ?, ?, ?)`
	_, err := db.Exec(query, n.UserId, n.ActorId, n.Type, n.PostId, n.CommentId, n.TimestampString)
	return err
}

func CreateNotification(notification Notification) error {
	if notification.ActorId == notification.UserId {
		return nil
	}
	notification.TimestampString = utils.GetCurrentTimestamp()
	return notification.Add()
}

func (user *User) MarkAsReadPost(post Post) error {
	for i, notification := range user.Notifications {
		if notification.PostId == post.Id && !notification.Read {
			err := user.MarkNotificationAsRead(notification.Id)
			if err != nil {
				return err
			}
			user.Notifications[i].Read = true
			user.UnreadNotificationsCount--
		}
	}
	return nil
}

func (comment *Comment) CreateCommentNotification(post Post) error {
	notification := Notification{
		UserId:    post.User.Id,
		ActorId:   comment.UserId,
		Type:      "comment",
		PostId:	   comment.PostId,
		CommentId: comment.Id,
	}
	err := CreateNotification(notification)
	if err != nil {
		return err
	}
	return nil
}

func (comment *Comment) CreateReactionNotification(userId int64, t string) error {
	notification := Notification{
		UserId:    comment.UserId,
		ActorId:   userId,
		CommentId: comment.Id,
		PostId:    comment.PostId,
		Type: t,
	}
	err := CreateNotification(notification)
	if err != nil {
		return err
	}
	return nil
}

func (post *Post) CreateReactionNotification(userId int64, t string) error {
	notification := Notification{
		UserId:  post.User.Id,
		ActorId: userId,
		Type: t,
		PostId: post.Id,
	}
	err := CreateNotification(notification)
	if err != nil {
		return err
	}
	return nil
}
