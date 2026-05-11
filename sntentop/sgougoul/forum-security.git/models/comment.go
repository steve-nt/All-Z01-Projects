package models

// Comment represents a comment made by a user on a post.
type Comment struct {
	ID        int
	PostID    int
	UserID    int
	Username  string
	Content   string
	CreatedAt string
}
