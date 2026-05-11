package db

import "time"

// ModeratorRequest represents one user request to become a moderator.
type ModeratorRequest struct {
	ID         int
	UserID     int
	Username   string
	Status     string
	CreatedAt  string
	ReviewedBy int
	ReviewedAt string
}

// CreateModeratorRequest creates a new moderator-role request for a normal user.
// AUDIT: users are allowed to request promotion, but only admins may approve it.
func CreateModeratorRequest(userID int) error {
	_, err := DB.Exec(`
		INSERT INTO moderator_requests (user_id, status)
		VALUES (?, 'pending')
	`, userID)
	return err
}

// HasPendingModeratorRequest reports whether the user already has a pending request.
// AUDIT: this prevents duplicate moderator requests from the same account.
func HasPendingModeratorRequest(userID int) (bool, error) {
	var count int
	err := DB.QueryRow(`
		SELECT COUNT(*)
		FROM moderator_requests
		WHERE user_id = ? AND status = 'pending'
	`, userID).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// GetAllModeratorRequests returns all moderator requests newest first.
// Includes the requesting username for the admin dashboard.
func GetAllModeratorRequests() ([]ModeratorRequest, error) {
	rows, err := DB.Query(`
		SELECT
			mr.id,
			mr.user_id,
			u.username,
			mr.status,
			mr.created_at,
			IFNULL(mr.reviewed_by, 0),
			IFNULL(mr.reviewed_at, '')
		FROM moderator_requests mr
		JOIN users u ON u.id = mr.user_id
		ORDER BY
			CASE WHEN mr.status = 'pending' THEN 0 ELSE 1 END,
			mr.created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []ModeratorRequest

	for rows.Next() {
		var mr ModeratorRequest
		if err := rows.Scan(
			&mr.ID,
			&mr.UserID,
			&mr.Username,
			&mr.Status,
			&mr.CreatedAt,
			&mr.ReviewedBy,
			&mr.ReviewedAt,
		); err != nil {
			return nil, err
		}
		out = append(out, mr)
	}

	return out, rows.Err()
}

// ApproveModeratorRequest approves a pending request and promotes the user.
// AUDIT: approval is done in a transaction so role update and request update stay consistent.
func ApproveModeratorRequest(requestID, adminUserID int) error {
	tx, err := DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var userID int
	err = tx.QueryRow(`
		SELECT user_id
		FROM moderator_requests
		WHERE id = ? AND status = 'pending'
	`, requestID).Scan(&userID)
	if err != nil {
		return err
	}

	if _, err := tx.Exec(`
		UPDATE users
		SET role = 'moderator'
		WHERE id = ?
	`, userID); err != nil {
		return err
	}

	if _, err := tx.Exec(`
		UPDATE moderator_requests
		SET status = 'approved',
		    reviewed_by = ?,
		    reviewed_at = ?
		WHERE id = ?
	`, adminUserID, time.Now().Format("2006-01-02 15:04:05"), requestID); err != nil {
		return err
	}

	return tx.Commit()
}

// RejectModeratorRequest rejects a pending moderator request.
func RejectModeratorRequest(requestID, adminUserID int) error {
	_, err := DB.Exec(`
		UPDATE moderator_requests
		SET status = 'rejected',
		    reviewed_by = ?,
		    reviewed_at = ?
		WHERE id = ? AND status = 'pending'
	`, adminUserID, time.Now().Format("2006-01-02 15:04:05"), requestID)
	return err
}