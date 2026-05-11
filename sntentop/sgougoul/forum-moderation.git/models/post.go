package models

// Post represents a forum post created by a user.
type Post struct {
	ID         int
	UserID     int
	Username   string
	Title      string
	Content    string
	CreatedAt  string
	Categories []string
	Status     string
}