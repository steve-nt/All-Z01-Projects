package db

import (
	"strings"

	"forum/models"
)

// nextAssignedRole determines the role for a newly created account.
// The very first account becomes admin so the forum always has one administrator.
// All later accounts default to normal user.
func nextAssignedRole() string {
	var count int
	if err := DB.QueryRow(`SELECT COUNT(*) FROM users`).Scan(&count); err != nil {
		return "user"
	}
	if count == 0 {
		return "admin"
	}
	return "user"
}

// GetUserRole returns the normalized role of the given user.
// Missing or unknown roles are treated as "user".
func GetUserRole(userID int) (string, error) {
	var role string
	if err := DB.QueryRow(`SELECT role FROM users WHERE id = ?`, userID).Scan(&role); err != nil {
		return "", err
	}

	role = strings.ToLower(strings.TrimSpace(role))
	switch role {
	case "admin", "moderator", "user":
		return role, nil
	default:
		return "user", nil
	}
}

// IsModeratorOrAdmin reports whether the user has elevated moderation access.
func IsModeratorOrAdmin(userID int) (bool, error) {
	role, err := GetUserRole(userID)
	if err != nil {
		return false, err
	}
	return role == "moderator" || role == "admin", nil
}

// IsAdmin reports whether the user is an administrator.
func IsAdmin(userID int) (bool, error) {
	role, err := GetUserRole(userID)
	if err != nil {
		return false, err
	}
	return role == "admin", nil
}

// GetAllUsers returns all users for admin role management.
// AUDIT: this supports admin promotion/demotion without changing existing auth flows.
func GetAllUsers() ([]models.User, error) {
	rows, err := DB.Query(`
		SELECT id, email, username, password,
		       IFNULL(provider, ''), IFNULL(provider_id, ''), role
		FROM users
		ORDER BY username ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []models.User

	for rows.Next() {
		var u models.User
		if err := rows.Scan(
			&u.ID,
			&u.Email,
			&u.Username,
			&u.Password,
			&u.Provider,
			&u.ProviderID,
			&u.Role,
		); err != nil {
			return nil, err
		}
		out = append(out, u)
	}

	return out, rows.Err()
}

// UpdateUserRole changes a user's role.
// AUDIT: role changes are reserved for admins in the handler layer.
func UpdateUserRole(userID int, role string) error {
	role = strings.ToLower(strings.TrimSpace(role))
	switch role {
	case "user", "moderator", "admin":
	default:
		role = "user"
	}

	_, err := DB.Exec(`UPDATE users SET role = ? WHERE id = ?`, role, userID)
	return err
}

// GetAdminUserIDs returns all administrator user IDs.
// AUDIT: used to fan out moderation-related notifications to admin accounts.
func GetAdminUserIDs() ([]int, error) {
	rows, err := DB.Query(`SELECT id FROM users WHERE role = 'admin' ORDER BY id ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []int
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		out = append(out, id)
	}

	return out, rows.Err()
}

// GetModeratorAndAdminUserIDs returns all moderator/admin IDs.
// AUDIT: new pending posts should notify every role that can review moderation.
func GetModeratorAndAdminUserIDs() ([]int, error) {
	rows, err := DB.Query(`
		SELECT id
		FROM users
		WHERE role IN ('moderator', 'admin')
		ORDER BY id ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []int
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		out = append(out, id)
	}

	return out, rows.Err()
}