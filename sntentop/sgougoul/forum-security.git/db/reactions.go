package db

import (
	"fmt"
	"strings"
)

// ReactionCounts holds the total likes/dislikes for a post or comment.
type ReactionCounts struct {
	Likes    int
	Dislikes int
}

// UpsertPostReaction inserts or updates a user's reaction for a post.
// value must be +1 (like) or -1 (dislike).
func UpsertPostReaction(userID, postID, value int) error {
	_, err := DB.Exec(`
		INSERT INTO reactions (user_id, post_id, comment_id, value)
		VALUES (?, ?, NULL, ?)
		ON CONFLICT(user_id, post_id)
		DO UPDATE SET value = excluded.value, created_at = CURRENT_TIMESTAMP
	`, userID, postID, value)
	return err
}

// UpsertCommentReaction inserts or updates a user's reaction for a comment.
// value must be +1 (like) or -1 (dislike).
func UpsertCommentReaction(userID, commentID, value int) error {
	_, err := DB.Exec(`
		INSERT INTO reactions (user_id, post_id, comment_id, value)
		VALUES (?, NULL, ?, ?)
		ON CONFLICT(user_id, comment_id)
		DO UPDATE SET value = excluded.value, created_at = CURRENT_TIMESTAMP
	`, userID, commentID, value)
	return err
}

// GetPostReactionCounts returns the like/dislike totals for one post.
func GetPostReactionCounts(postID int) (ReactionCounts, error) {
	var likes, dislikes int

	// Aggregate counts directly in SQL for correctness + speed.
	err := DB.QueryRow(`
		SELECT
			IFNULL(SUM(CASE WHEN value = 1 THEN 1 ELSE 0 END), 0) AS likes,
			IFNULL(SUM(CASE WHEN value = -1 THEN 1 ELSE 0 END), 0) AS dislikes
		FROM reactions
		WHERE post_id = ?
	`, postID).Scan(&likes, &dislikes)

	return ReactionCounts{Likes: likes, Dislikes: dislikes}, err
}

// GetCommentCountsByPost returns a map: comment_id -> {likes, dislikes}
// for every comment under a specific post.
func GetCommentCountsByPost(postID int) (map[int]ReactionCounts, error) {
	rows, err := DB.Query(`
		SELECT c.id,
			IFNULL(SUM(CASE WHEN r.value = 1 THEN 1 ELSE 0 END), 0) AS likes,
			IFNULL(SUM(CASE WHEN r.value = -1 THEN 1 ELSE 0 END), 0) AS dislikes
		FROM comments c
		LEFT JOIN reactions r ON r.comment_id = c.id
		WHERE c.post_id = ?
		GROUP BY c.id
	`, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make(map[int]ReactionCounts)

	for rows.Next() {
		var id, likes, dislikes int
		if err := rows.Scan(&id, &likes, &dislikes); err != nil {
			return nil, err
		}
		out[id] = ReactionCounts{Likes: likes, Dislikes: dislikes}
	}

	return out, rows.Err()
}

// GetPostCountsByPostIDs fetches reaction totals for many posts in one query.
// Returns a map: post_id -> {likes, dislikes}.
func GetPostCountsByPostIDs(postIDs []int) (map[int]ReactionCounts, error) {
	out := make(map[int]ReactionCounts)

	// No posts -> return empty map without querying DB.
	if len(postIDs) == 0 {
		return out, nil
	}

	// Build "?, ?, ?, ..." placeholders for the IN (...) query.
	placeholders := strings.TrimRight(strings.Repeat("?,", len(postIDs)), ",")

	// Convert []int into []any for DB.Query args.
	args := make([]any, 0, len(postIDs))
	for _, id := range postIDs {
		args = append(args, id)
	}

	// Use fmt.Sprintf only for placeholders; values remain parameterized (safe).
	q := fmt.Sprintf(`
		SELECT post_id,
			IFNULL(SUM(CASE WHEN value = 1 THEN 1 ELSE 0 END), 0) AS likes,
			IFNULL(SUM(CASE WHEN value = -1 THEN 1 ELSE 0 END), 0) AS dislikes
		FROM reactions
		WHERE post_id IN (%s)
		GROUP BY post_id
	`, placeholders)

	rows, err := DB.Query(q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var pid, likes, dislikes int
		if err := rows.Scan(&pid, &likes, &dislikes); err != nil {
			return nil, err
		}
		out[pid] = ReactionCounts{Likes: likes, Dislikes: dislikes}
	}

	return out, rows.Err()
}
