package db

// ActivityReaction represents one post the user reacted to.
type ActivityReaction struct {
	PostID      int
	PostTitle   string
	Reaction    int
	CreatedAt   string
}

// ActivityComment represents one comment the user wrote,
// along with the related post information.
type ActivityComment struct {
	CommentID    int
	PostID       int
	PostTitle    string
	Content      string
	CreatedAt    string
}

// GetReactionsByUser returns post reactions made by a user.
// Only post reactions are included here because the activity page
// is focused on where the user reacted in discussions.
func GetReactionsByUser(userID int) ([]ActivityReaction, error) {
	rows, err := DB.Query(`
		SELECT
			p.id,
			p.title,
			r.value,
			r.created_at
		FROM reactions r
		JOIN posts p ON p.id = r.post_id
		WHERE r.user_id = ? AND r.post_id IS NOT NULL
		ORDER BY r.created_at DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []ActivityReaction

	for rows.Next() {
		var a ActivityReaction
		if err := rows.Scan(&a.PostID, &a.PostTitle, &a.Reaction, &a.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, a)
	}

	return out, rows.Err()
}

// GetCommentsByUser returns the comments written by the user,
// along with the title of the related post.
func GetCommentsByUser(userID int) ([]ActivityComment, error) {
	rows, err := DB.Query(`
		SELECT
			c.id,
			c.post_id,
			p.title,
			c.content,
			c.created_at
		FROM comments c
		JOIN posts p ON p.id = c.post_id
		WHERE c.user_id = ?
		ORDER BY c.created_at DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []ActivityComment

	for rows.Next() {
		var a ActivityComment
		if err := rows.Scan(&a.CommentID, &a.PostID, &a.PostTitle, &a.Content, &a.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, a)
	}

	return out, rows.Err()
}