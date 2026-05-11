package db

import "time"

// Report represents one moderator report sent to administrators.
type Report struct {
	ID            int
	PostID        int
	PostTitle     string
	ReporterID    int
	ReporterName  string
	Reason        string
	Status        string
	AdminResponse string
	CreatedAt     string
	ReviewedBy    int
	ReviewedAt    string
}

// CreateReport stores a moderator report for an existing post.
// AUDIT: moderators use this flow to escalate suspicious or problematic content to admins.
func CreateReport(postID, reporterID int, reason string) error {
	_, err := DB.Exec(`
		INSERT INTO reports (post_id, reporter_id, reason, status)
		VALUES (?, ?, ?, 'pending')
	`, postID, reporterID, reason)
	return err
}

// GetAllReports returns all reports for the admin dashboard.
// Includes the post title and reporter username for easier review.
func GetAllReports() ([]Report, error) {
	rows, err := DB.Query(`
		SELECT
			r.id,
			r.post_id,
			p.title,
			r.reporter_id,
			u.username,
			r.reason,
			r.status,
			IFNULL(r.admin_response, ''),
			r.created_at,
			IFNULL(r.reviewed_by, 0),
			IFNULL(r.reviewed_at, '')
		FROM reports r
		JOIN posts p ON p.id = r.post_id
		JOIN users u ON u.id = r.reporter_id
		ORDER BY
			CASE WHEN r.status = 'pending' THEN 0 ELSE 1 END,
			r.created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []Report

	for rows.Next() {
		var rep Report
		if err := rows.Scan(
			&rep.ID,
			&rep.PostID,
			&rep.PostTitle,
			&rep.ReporterID,
			&rep.ReporterName,
			&rep.Reason,
			&rep.Status,
			&rep.AdminResponse,
			&rep.CreatedAt,
			&rep.ReviewedBy,
			&rep.ReviewedAt,
		); err != nil {
			return nil, err
		}
		out = append(out, rep)
	}

	return out, rows.Err()
}

// ResolveReport marks a report as resolved and stores the admin response.
// AUDIT: this gives administrators a traceable response to moderator escalation.
func ResolveReport(reportID, adminUserID int, response string) error {
	_, err := DB.Exec(`
		UPDATE reports
		SET status = 'resolved',
		    admin_response = ?,
		    reviewed_by = ?,
		    reviewed_at = ?
		WHERE id = ?
	`, response, adminUserID, time.Now().Format("2006-01-02 15:04:05"), reportID)
	return err
}

// HasPendingReportForPost checks whether the same post already has a pending report.
// AUDIT: this avoids noisy duplicate escalation for the same unresolved post.
func HasPendingReportForPost(postID int) (bool, error) {
	var count int
	err := DB.QueryRow(`
		SELECT COUNT(*)
		FROM reports
		WHERE post_id = ? AND status = 'pending'
	`, postID).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}