package db

// migrateUsersOAuthColumns ensures OAuth columns exist on the users table.
// This allows existing databases (created before OAuth support)
// to be upgraded safely without losing data.
func migrateUsersOAuthColumns() error {
	cols, err := getTableColumns("users")
	if err != nil {
		return err
	}

	// Add provider column if missing
	if !cols["provider"] {
		if _, err := DB.Exec(`ALTER TABLE users ADD COLUMN provider TEXT;`); err != nil {
			return err
		}
	}

	// Add provider_id column if missing
	if !cols["provider_id"] {
		if _, err := DB.Exec(`ALTER TABLE users ADD COLUMN provider_id TEXT;`); err != nil {
			return err
		}
	}

	return nil
}

// migrateUserRoles ensures the role column exists on the users table.
// Existing users are initialized as normal users first.
// Administrative bootstrap is handled separately by ensureAdminExists.
func migrateUserRoles() error {
	cols, err := getTableColumns("users")
	if err != nil {
		return err
	}

	if !cols["role"] {
		if _, err := DB.Exec(`ALTER TABLE users ADD COLUMN role TEXT NOT NULL DEFAULT 'user';`); err != nil {
			return err
		}
	}

	_, err = DB.Exec(`UPDATE users SET role = 'user' WHERE role IS NULL OR TRIM(role) = '';`)
	return err
}

// migratePostModerationColumns ensures moderation columns exist on posts.
// Existing posts are marked approved so older forum data stays publicly visible.
func migratePostModerationColumns() error {
	cols, err := getTableColumns("posts")
	if err != nil {
		return err
	}

	if !cols["status"] {
		if _, err := DB.Exec(`ALTER TABLE posts ADD COLUMN status TEXT NOT NULL DEFAULT 'approved';`); err != nil {
			return err
		}
	}

	if !cols["reviewed_by"] {
		if _, err := DB.Exec(`ALTER TABLE posts ADD COLUMN reviewed_by INTEGER;`); err != nil {
			return err
		}
	}

	if !cols["reviewed_at"] {
		if _, err := DB.Exec(`ALTER TABLE posts ADD COLUMN reviewed_at DATETIME;`); err != nil {
			return err
		}
	}

	_, err = DB.Exec(`UPDATE posts SET status = 'approved' WHERE status IS NULL OR TRIM(status) = '';`)
	return err
}

// migrateModeratorRequestsTable ensures the moderator request table exists.
// AUDIT: older databases may not have the moderation-request workflow yet.
func migrateModeratorRequestsTable() error {
	_, err := DB.Exec(`
		CREATE TABLE IF NOT EXISTS moderator_requests (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			status TEXT NOT NULL DEFAULT 'pending',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			reviewed_by INTEGER,
			reviewed_at DATETIME,
			FOREIGN KEY (user_id) REFERENCES users(id),
			FOREIGN KEY (reviewed_by) REFERENCES users(id)
		);
	`)
	return err
}

// migrateReportsTable ensures the reports table exists.
// AUDIT: older databases may not have moderator -> admin reporting yet.
func migrateReportsTable() error {
	_, err := DB.Exec(`
		CREATE TABLE IF NOT EXISTS reports (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			post_id INTEGER NOT NULL,
			reporter_id INTEGER NOT NULL,
			reason TEXT NOT NULL,
			status TEXT NOT NULL DEFAULT 'pending',
			admin_response TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			reviewed_by INTEGER,
			reviewed_at DATETIME,
			FOREIGN KEY (post_id) REFERENCES posts(id),
			FOREIGN KEY (reporter_id) REFERENCES users(id),
			FOREIGN KEY (reviewed_by) REFERENCES users(id)
		);
	`)
	return err
}

// migrateNotificationsForModeration ensures notifications can support
// generic moderation messages that are not always tied to a post link.
func migrateNotificationsForModeration() error {
	cols, err := getTableColumns("notifications")
	if err != nil {
		return err
	}

	// Older databases may not have a free-text message column yet.
	if !cols["message"] {
		if _, err := DB.Exec(`ALTER TABLE notifications ADD COLUMN message TEXT NOT NULL DEFAULT '';`); err != nil {
			return err
		}
	}

	return nil
}