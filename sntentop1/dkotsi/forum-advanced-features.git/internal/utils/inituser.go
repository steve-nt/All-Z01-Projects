package utils

import "forum-advanced-features/internal/backend/models"

func InitUser() models.User {

	comment := []models.Comment{}

	post := ConstructPost(comment)
	likedposts := []models.Post{}
	likedposts = append(likedposts, post)
	notification := ConstructNotification(post)
	slnotification := []models.Notification{}
	slnotification = append(slnotification, notification)
	user := ConstructUser(likedposts, slnotification)
	return user
}

func ConstructPost(comment []models.Comment) models.Post {
	return models.Post{
		Comments: comment,
	}
}

func ConstructNotification(post models.Post) models.Notification {
	return models.Notification{
		Post: post,
	}
}

func ConstructUser(likedpost []models.Post, slnotif []models.Notification) models.User {
	return models.User{
		Notifications: slnotif,
		LikedPosts:    likedpost,
		CreatedPosts:  likedpost,
	}
}
