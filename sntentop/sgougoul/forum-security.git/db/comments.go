package db

import "forum/models"

// CreateComment inserts a new comment into the database.
// postID  -> the post this comment belongs to
// userID  -> the author of the comment
// content -> the comment text
func CreateComment(postID, userID int, content string) error {
	_, err := DB.Exec(
		`INSERT INTO comments (post_id, user_id, content)
		 VALUES (?, ?, ?)`,
		postID, userID, content,
	)
	return err
}

// GetCommentByID returns one comment by its ID.
// Used for edit/delete ownership checks.
func GetCommentByID(commentID int) (models.Comment, error) {
	var c models.Comment

	err := DB.QueryRow(
		`SELECT id, post_id, user_id, content, created_at
		 FROM comments
		 WHERE id = ?`,
		commentID,
	).Scan(&c.ID, &c.PostID, &c.UserID, &c.Content, &c.CreatedAt)

	return c, err
}

// UpdateComment changes the content of an existing comment.
func UpdateComment(commentID int, content string) error {
	_, err := DB.Exec(
		`UPDATE comments
		 SET content = ?
		 WHERE id = ?`,
		content, commentID,
	)
	return err
}

// DeleteComment removes a comment and related comment reactions.
func DeleteComment(commentID int) error {
	_, err := DB.Exec(`
		DELETE FROM reactions
		WHERE comment_id = ?
	`, commentID)
	if err != nil {
		return err
	}

	_, err = DB.Exec(`DELETE FROM comments WHERE id = ?`, commentID)
	return err
}

// GetCommentsByPostID retrieves all comments for a given post.
//
// Includes the username by joining the users table.
// Comments are returned oldest → newest so discussions read naturally.
func GetCommentsByPostID(postID int) ([]models.Comment, error) {
	rows, err := DB.Query(
		`SELECT c.id, c.post_id, c.user_id, u.username, c.content, c.created_at
		 FROM comments c
		 JOIN users u ON u.id = c.user_id
		 WHERE c.post_id = ?
		 ORDER BY c.created_at ASC`,
		postID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []models.Comment

	for rows.Next() {
		var c models.Comment

		// Scan DB row into Comment struct
		if err := rows.Scan(
			&c.ID,
			&c.PostID,
			&c.UserID,
			&c.Username,
			&c.Content,
			&c.CreatedAt,
		); err != nil {
			return nil, err
		}

		comments = append(comments, c)
	}

	return comments, rows.Err()
}