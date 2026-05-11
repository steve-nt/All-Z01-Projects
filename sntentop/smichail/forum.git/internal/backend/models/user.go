package models

import "time"

type User struct {
	UUID          string
	Mail          string
	Username      string
	Password      string
	Role          string
	Verified      bool
	CreationDate  time.Time
	Notifications []Notification
	Activities    []Activity
	LikedPosts    []Post
	CreatedPosts  []Post
}

type Notification struct {
	ID       int
	Username string
	Action   string
	Post     Post
	Comment  Comment
	Seen     bool
}

type Activity struct {
	Action       string
	Post         Post
	Comment      Comment
	CreationDate string
	UUID         string
}
