package db

import "log"

// createTables creates all required tables for the forum.
// Uses IF NOT EXISTS so it is safe to run every startup.
func createTables() {
	queries := []string{
		// Users table
		// Stores authentication credentials and forum role.
		`CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			email TEXT UNIQUE NOT NULL,
			username TEXT NOT NULL,
			password TEXT NOT NULL,
			role TEXT NOT NULL DEFAULT 'user'
		);`,

		// Sessions table
		// Each user can only have one active session (user_id UNIQUE).
		`CREATE TABLE IF NOT EXISTS sessions (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER UNIQUE NOT NULL,
			token TEXT UNIQUE NOT NULL,
			expires_at DATETIME NOT NULL,
			FOREIGN KEY (user_id) REFERENCES users(id)
		);`,

		// Posts table
		// status controls moderation visibility:
		// - pending  -> awaiting approval
		// - approved -> publicly visible
		// - rejected -> hidden from public lists
		`CREATE TABLE IF NOT EXISTS posts (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			title TEXT NOT NULL,
			content TEXT NOT NULL,
			status TEXT NOT NULL DEFAULT 'pending',
			reviewed_by INTEGER,
			reviewed_at DATETIME,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id)
		);`,

		// Comments table
		// Stores comments attached to posts.
		`CREATE TABLE IF NOT EXISTS comments (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			post_id INTEGER NOT NULL,
			user_id INTEGER NOT NULL,
			content TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (post_id) REFERENCES posts(id),
			FOREIGN KEY (user_id) REFERENCES users(id)
		);`,

		// Categories table
		// Stores predefined forum categories.
		`CREATE TABLE IF NOT EXISTS categories (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT UNIQUE NOT NULL
		);`,

		// Post-Category relation table (many-to-many)
		`CREATE TABLE IF NOT EXISTS post_categories (
			post_id INTEGER NOT NULL,
			category_id INTEGER NOT NULL,
			PRIMARY KEY (post_id, category_id),
			FOREIGN KEY (post_id) REFERENCES posts(id),
			FOREIGN KEY (category_id) REFERENCES categories(id)
		);`,

		// Reactions table
		// Stores likes/dislikes for BOTH posts and comments.
		// Exactly one of post_id OR comment_id must be set.
		`CREATE TABLE IF NOT EXISTS reactions (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			post_id INTEGER,
			comment_id INTEGER,
			value INTEGER NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,

			FOREIGN KEY (user_id) REFERENCES users(id),
			FOREIGN KEY (post_id) REFERENCES posts(id),
			FOREIGN KEY (comment_id) REFERENCES comments(id),

			-- Only one reaction per user per post/comment
			UNIQUE(user_id, post_id),
			UNIQUE(user_id, comment_id),

			-- Must be either post OR comment reaction
			CHECK (
				(post_id IS NOT NULL AND comment_id IS NULL) OR
				(post_id IS NULL AND comment_id IS NOT NULL)
			),

			-- Only like (+1) or dislike (-1)
			CHECK (value IN (1, -1))
		);`,

		// Notifications table
		// Stores user notifications related to their posts and moderation events.
		`CREATE TABLE IF NOT EXISTS notifications (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			actor_user_id INTEGER NOT NULL,
			post_id INTEGER,
			type TEXT NOT NULL,
			message TEXT NOT NULL DEFAULT '',
			is_read INTEGER NOT NULL DEFAULT 0,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,

			FOREIGN KEY (user_id) REFERENCES users(id),
			FOREIGN KEY (actor_user_id) REFERENCES users(id),
			FOREIGN KEY (post_id) REFERENCES posts(id)
		);`,

		// Moderator requests table
		// AUDIT: normal users can request promotion to moderator.
		// Admin reviews these requests from the moderation dashboard.
		`CREATE TABLE IF NOT EXISTS moderator_requests (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			status TEXT NOT NULL DEFAULT 'pending',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			reviewed_by INTEGER,
			reviewed_at DATETIME,
			FOREIGN KEY (user_id) REFERENCES users(id),
			FOREIGN KEY (reviewed_by) REFERENCES users(id)
		);`,

		// Reports table
		// AUDIT: moderators can report problematic posts to administrators.
		// Admin can later resolve these reports.
		`CREATE TABLE IF NOT EXISTS reports (
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
		);`,
	}

	// Execute all table creation queries
	for _, q := range queries {
		if _, err := DB.Exec(q); err != nil {
			log.Fatal(err)
		}
	}
}

// createIndexes adds additional database constraints.
func createIndexes() {
	// Unique username requirement
	if _, err := DB.Exec(`CREATE UNIQUE INDEX IF NOT EXISTS idx_users_username ON users(username);`); err != nil {
		log.Println("WARNING: could not create unique username index:", err)
	}

	// Ensure uniqueness of OAuth identities (provider + provider_id)
	// Multiple NULL values are allowed for normal password users.
	if _, err := DB.Exec(`CREATE UNIQUE INDEX IF NOT EXISTS idx_users_provider_identity ON users(provider, provider_id);`); err != nil {
		log.Println("WARNING: could not create provider identity index:", err)
	}

	// Notification indexes improve loading unread notifications and activity lists.
	_, _ = DB.Exec(`CREATE INDEX IF NOT EXISTS idx_notifications_user_id ON notifications(user_id);`)
	_, _ = DB.Exec(`CREATE INDEX IF NOT EXISTS idx_notifications_is_read ON notifications(is_read);`)
	_, _ = DB.Exec(`CREATE INDEX IF NOT EXISTS idx_notifications_created_at ON notifications(created_at);`)

	// Moderation indexes improve filtering by role and post status.
	_, _ = DB.Exec(`CREATE INDEX IF NOT EXISTS idx_users_role ON users(role);`)
	_, _ = DB.Exec(`CREATE INDEX IF NOT EXISTS idx_posts_status ON posts(status);`)

	// AUDIT: indexes for moderation request workflow.
	_, _ = DB.Exec(`CREATE INDEX IF NOT EXISTS idx_moderator_requests_user_id ON moderator_requests(user_id);`)
	_, _ = DB.Exec(`CREATE INDEX IF NOT EXISTS idx_moderator_requests_status ON moderator_requests(status);`)

	// AUDIT: indexes for moderator reports and admin review.
	_, _ = DB.Exec(`CREATE INDEX IF NOT EXISTS idx_reports_post_id ON reports(post_id);`)
	_, _ = DB.Exec(`CREATE INDEX IF NOT EXISTS idx_reports_reporter_id ON reports(reporter_id);`)
	_, _ = DB.Exec(`CREATE INDEX IF NOT EXISTS idx_reports_status ON reports(status);`)
}

// seedCategories inserts default categories on first run.
// INSERT OR IGNORE prevents duplicates.
func seedCategories() {
	defaultCategories := []string{
		"General",
		"Technology",
		"Gaming",
		"Movies & TV",
		"Music",
		"Sports",
		"Help",
	}

	for _, name := range defaultCategories {
		if _, err := DB.Exec(`INSERT OR IGNORE INTO categories (name) VALUES (?)`, name); err != nil {
			log.Println("seedCategories error:", err)
		}
	}
}