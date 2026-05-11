package db

import (
	"database/sql"
	"strings"

	"forum/models"
)

// GetVisiblePostByID returns a post if the current viewer is allowed to see it.
// AUDIT:
// - approved posts are public
// - owners may view their own posts even when pending/rejected
// - moderators/admins may review any post
func GetVisiblePostByID(postID, viewerUserID int, isPrivileged bool) (models.Post, error) {
	post, err := GetPostByIDWithCategories(postID)
	if err != nil {
		return post, err
	}

	var status string
	if err := DB.QueryRow(`SELECT status FROM posts WHERE id = ?`, postID).Scan(&status); err != nil {
		return post, err
	}
	post.Status = strings.TrimSpace(status)

	if post.Status == "approved" {
		return post, nil
	}

	if post.UserID == viewerUserID {
		return post, nil
	}

	if isPrivileged {
		return post, nil
	}

	return models.Post{}, sql.ErrNoRows
}

// GetPendingPosts returns all posts waiting for moderator/admin review.
// AUDIT: this is the moderation queue used by privileged roles.
func GetPendingPosts() ([]models.Post, error) {
	rows, err := DB.Query(`
		SELECT p.id, p.user_id, u.username, p.title, p.content, p.created_at,
		       IFNULL(GROUP_CONCAT(c.name, ','), '') AS cats,
		       p.status
		FROM posts p
		JOIN users u ON u.id = p.user_id
		LEFT JOIN post_categories pc ON pc.post_id = p.id
		LEFT JOIN categories c ON c.id = pc.category_id
		WHERE p.status = 'pending'
		GROUP BY p.id
		ORDER BY p.created_at ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []models.Post

	for rows.Next() {
		var p models.Post
		var cats string

		if err := rows.Scan(&p.ID, &p.UserID, &p.Username, &p.Title, &p.Content, &p.CreatedAt, &cats, &p.Status); err != nil {
			return nil, err
		}

		p.Categories = parseCategoriesCSV(cats)
		posts = append(posts, p)
	}

	return posts, rows.Err()
}

// ApprovePost marks a pending post as approved.
// AUDIT: moderators/admins can publish content after review.
func ApprovePost(postID, reviewerUserID int) error {
	_, err := DB.Exec(`
		UPDATE posts
		SET status = 'approved',
		    reviewed_by = ?,
		    reviewed_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`, reviewerUserID, postID)
	return err
}

// RejectPost marks a pending post as rejected.
// AUDIT: rejected posts remain hidden from public lists.
func RejectPost(postID, reviewerUserID int) error {
	_, err := DB.Exec(`
		UPDATE posts
		SET status = 'rejected',
		    reviewed_by = ?,
		    reviewed_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`, reviewerUserID, postID)
	return err
}

// ResetPostToPending moves an edited post back into the moderation queue.
// AUDIT: when a user edits a post after approval/rejection, moderators/admins
// should review the updated content before it becomes public again.
func ResetPostToPending(postID int) error {
	_, err := DB.Exec(`
		UPDATE posts
		SET status = 'pending',
		    reviewed_by = NULL,
		    reviewed_at = NULL
		WHERE id = ?
	`, postID)
	return err
}