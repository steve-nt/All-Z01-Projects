package db

import (
	"database/sql"
	"strings"
)

// ensureAdminExists guarantees that at least one admin user exists.
// If the database already contains users but no admin yet, the oldest user becomes admin.
func ensureAdminExists() error {
	var adminCount int
	if err := DB.QueryRow(`SELECT COUNT(*) FROM users WHERE role = 'admin'`).Scan(&adminCount); err != nil {
		return err
	}
	if adminCount > 0 {
		return nil
	}

	var firstUserID int
	err := DB.QueryRow(`SELECT id FROM users ORDER BY id ASC LIMIT 1`).Scan(&firstUserID)
	if err == sql.ErrNoRows {
		return nil
	}
	if err != nil {
		return err
	}

	_, err = DB.Exec(`UPDATE users SET role = 'admin' WHERE id = ?`, firstUserID)
	return err
}

// getTableColumns reads existing column names from a table.
// Used by migrations to detect whether a column already exists.
func getTableColumns(table string) (map[string]bool, error) {
	rows, err := DB.Query(`PRAGMA table_info(` + table + `);`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make(map[string]bool)

	for rows.Next() {
		var (
			cid       int
			name      string
			typ       string
			notnull   int
			dfltValue sql.NullString
			pk        int
		)

		if err := rows.Scan(&cid, &name, &typ, &notnull, &dfltValue, &pk); err != nil {
			return nil, err
		}

		out[strings.ToLower(strings.TrimSpace(name))] = true
	}

	return out, rows.Err()
}