package db

import (
	"database/sql"
	"strings"

	"forum/models"
)

// CreatePost inserts a new post in the `posts` table and returns the new post ID.
// AUDIT: status is set explicitly to "pending" so migrated databases cannot
// accidentally auto-approve new posts due to older column defaults.
func CreatePost(userID int, title, content string) (int, error) {
	res, err := DB.Exec(
		`INSERT INTO posts (user_id, title, content, status) VALUES (?, ?, ?, 'pending')`,
		userID, title, content,
	)
	if err != nil {
		return 0, err
	}

	id64, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id64), nil
}

// UpdatePost updates the title and content of an existing post.
// Ownership checks are handled in the handler layer before calling this.
func UpdatePost(postID int, title, content string) error {
	_, err := DB.Exec(
		`UPDATE posts
		 SET title = ?, content = ?
		 WHERE id = ?`,
		title, content, postID,
	)
	return err
}

// DeletePost removes a post and all related data.
// This includes post-category links, notifications, reactions and comments.
func DeletePost(postID int) error {
	tx, err := DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Remove comment reactions first (they depend on comments).
	if _, err := tx.Exec(`
		DELETE FROM reactions
		WHERE comment_id IN (
			SELECT id FROM comments WHERE post_id = ?
		)
	`, postID); err != nil {
		return err
	}

	// Remove comments under the post.
	if _, err := tx.Exec(`DELETE FROM comments WHERE post_id = ?`, postID); err != nil {
		return err
	}

	// Remove direct post reactions.
	if _, err := tx.Exec(`DELETE FROM reactions WHERE post_id = ?`, postID); err != nil {
		return err
	}

	// Remove post-category relations.
	if _, err := tx.Exec(`DELETE FROM post_categories WHERE post_id = ?`, postID); err != nil {
		return err
	}

	// Remove notifications related to the post.
	if _, err := tx.Exec(`DELETE FROM notifications WHERE post_id = ?`, postID); err != nil {
		return err
	}

	// AUDIT: reports are part of the moderation workflow and must be removed with the post.
	if _, err := tx.Exec(`DELETE FROM reports WHERE post_id = ?`, postID); err != nil {
		return err
	}

	// Remove the post itself.
	if _, err := tx.Exec(`DELETE FROM posts WHERE id = ?`, postID); err != nil {
		return err
	}

	return tx.Commit()
}

// GetPostByID returns a single post by ID (without categories).
func GetPostByID(id int) (models.Post, error) {
	var p models.Post

	err := DB.QueryRow(
		`SELECT p.id, p.user_id, u.username, p.title, p.content, p.created_at
		 FROM posts p
		 JOIN users u ON u.id = p.user_id
		 WHERE p.id = ?`,
		id,
	).Scan(&p.ID, &p.UserID, &p.Username, &p.Title, &p.Content, &p.CreatedAt)

	return p, err
}

// GetPostOwnerID returns the author user_id for a given post.
// Used for ownership checks and notifications.
func GetPostOwnerID(postID int) (int, error) {
	var userID int
	err := DB.QueryRow(`SELECT user_id FROM posts WHERE id = ?`, postID).Scan(&userID)
	return userID, err
}

// GetPostByIDWithCategories returns a single post by ID and also loads its categories.
// Categories are aggregated via GROUP_CONCAT from the many-to-many relation.
func GetPostByIDWithCategories(id int) (models.Post, error) {
	var p models.Post
	var cats string

	err := DB.QueryRow(`
		SELECT p.id, p.user_id, u.username, p.title, p.content, p.created_at,
		       IFNULL(GROUP_CONCAT(c.name, ','), '') AS cats
		FROM posts p
		JOIN users u ON u.id = p.user_id
		LEFT JOIN post_categories pc ON pc.post_id = p.id
		LEFT JOIN categories c ON c.id = pc.category_id
		WHERE p.id = ?
		GROUP BY p.id
	`, id).Scan(&p.ID, &p.UserID, &p.Username, &p.Title, &p.Content, &p.CreatedAt, &cats)
	if err != nil {
		return p, err
	}

	p.Categories = parseCategoriesCSV(cats)
	return p, nil
}

// PostExists checks whether a post exists (used for validation before writing comments/reactions).
func PostExists(postID int) (bool, error) {
	var x int

	err := DB.QueryRow(`SELECT 1 FROM posts WHERE id = ?`, postID).Scan(&x)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return true, nil
}

// parseCategoriesCSV converts the GROUP_CONCAT string into a []string.
// Example: "Gaming,Technology" -> []string{"Gaming","Technology"}.
func parseCategoriesCSV(cats string) []string {
	cats = strings.TrimSpace(cats)
	if cats == "" {
		return nil
	}

	parts := strings.Split(cats, ",")
	out := make([]string, 0, len(parts))

	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}

	return out
}