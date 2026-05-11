package db

import (
	"strings"

	"forum/models"
)

// GetAllPosts returns all publicly visible posts and their categories.
// If category != "", the result is filtered to posts that have that category name.
// AUDIT: guests and normal users should only see approved posts.
func GetAllPosts(category string) ([]models.Post, error) {
	query := `
	SELECT p.id, p.user_id, u.username, p.title, p.content, p.created_at,
	       IFNULL(GROUP_CONCAT(c.name, ','), '') AS cats
	FROM posts p
	JOIN users u ON u.id = p.user_id
	LEFT JOIN post_categories pc ON pc.post_id = p.id
	LEFT JOIN categories c ON c.id = pc.category_id
	WHERE p.status = 'approved'
	`
	args := []any{}

	// Optional category filter (by category name)
	if category != "" {
		query += `
		AND p.id IN (
			SELECT pc2.post_id
			FROM post_categories pc2
			JOIN categories c2 ON c2.id = pc2.category_id
			WHERE c2.name = ?
		)
		`
		args = append(args, strings.TrimSpace(category))
	}

	query += `
	GROUP BY p.id
	ORDER BY p.created_at DESC
	`

	rows, err := DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []models.Post

	for rows.Next() {
		var p models.Post
		var cats string

		if err := rows.Scan(&p.ID, &p.UserID, &p.Username, &p.Title, &p.Content, &p.CreatedAt, &cats); err != nil {
			return nil, err
		}

		p.Categories = parseCategoriesCSV(cats)
		posts = append(posts, p)
	}

	return posts, rows.Err()
}

// GetRecentPosts returns the latest approved posts for the home page preview.
// `limit` defaults to 5 if a non-positive value is provided.
func GetRecentPosts(limit int) ([]models.Post, error) {
	if limit <= 0 {
		limit = 5
	}

	rows, err := DB.Query(`
		SELECT p.id, p.user_id, u.username, p.title, p.content, p.created_at,
		       IFNULL(GROUP_CONCAT(c.name, ','), '') AS cats
		FROM posts p
		JOIN users u ON u.id = p.user_id
		LEFT JOIN post_categories pc ON pc.post_id = p.id
		LEFT JOIN categories c ON c.id = pc.category_id
		WHERE p.status = 'approved'
		GROUP BY p.id
		ORDER BY p.created_at DESC
		LIMIT ?
	`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []models.Post

	for rows.Next() {
		var p models.Post
		var cats string

		if err := rows.Scan(&p.ID, &p.UserID, &p.Username, &p.Title, &p.Content, &p.CreatedAt, &cats); err != nil {
			return nil, err
		}

		p.Categories = parseCategoriesCSV(cats)
		posts = append(posts, p)
	}

	return posts, rows.Err()
}

// GetTopLikedPosts returns the most liked approved posts for the home page.
// Posts are ordered by total likes descending, then by newest first.
func GetTopLikedPosts(limit int) ([]models.Post, error) {
	if limit <= 0 {
		limit = 5
	}

	rows, err := DB.Query(`
		SELECT
			p.id,
			p.user_id,
			u.username,
			p.title,
			p.content,
			p.created_at,
			IFNULL(GROUP_CONCAT(DISTINCT c.name), '') AS cats
		FROM posts p
		JOIN users u ON u.id = p.user_id
		LEFT JOIN post_categories pc ON pc.post_id = p.id
		LEFT JOIN categories c ON c.id = pc.category_id
		LEFT JOIN reactions r ON r.post_id = p.id
		WHERE p.status = 'approved'
		GROUP BY p.id
		ORDER BY
			IFNULL(SUM(CASE WHEN r.value = 1 THEN 1 ELSE 0 END), 0) DESC,
			p.created_at DESC
		LIMIT ?
	`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []models.Post

	for rows.Next() {
		var p models.Post
		var cats string

		if err := rows.Scan(&p.ID, &p.UserID, &p.Username, &p.Title, &p.Content, &p.CreatedAt, &cats); err != nil {
			return nil, err
		}

		p.Categories = parseCategoriesCSV(cats)
		posts = append(posts, p)
	}

	return posts, rows.Err()
}

// GetPostsByUser returns all posts created by a specific user.
// AUDIT: users should still be able to see their own pending/rejected posts in "My Posts".
func GetPostsByUser(userID int) ([]models.Post, error) {
	rows, err := DB.Query(`
		SELECT p.id, p.user_id, u.username, p.title, p.content, p.created_at,
		       IFNULL(GROUP_CONCAT(c.name, ','), '') AS cats,
		       p.status
		FROM posts p
		JOIN users u ON u.id = p.user_id
		LEFT JOIN post_categories pc ON pc.post_id = p.id
		LEFT JOIN categories c ON c.id = pc.category_id
		WHERE p.user_id = ?
		GROUP BY p.id
		ORDER BY p.created_at DESC
	`, userID)
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

// GetLikedPostsByUser returns approved posts liked by a given user.
// A "like" is stored as reactions.value = 1 for that post.
func GetLikedPostsByUser(userID int) ([]models.Post, error) {
	rows, err := DB.Query(`
		SELECT p.id, p.user_id, u.username, p.title, p.content, p.created_at,
		       IFNULL(GROUP_CONCAT(c.name, ','), '') AS cats,
		       p.status
		FROM posts p
		JOIN users u ON u.id = p.user_id
		JOIN reactions r ON r.post_id = p.id
		LEFT JOIN post_categories pc ON pc.post_id = p.id
		LEFT JOIN categories c ON c.id = pc.category_id
		WHERE r.user_id = ? AND r.value = 1 AND p.status = 'approved'
		GROUP BY p.id
		ORDER BY p.created_at DESC
	`, userID)
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