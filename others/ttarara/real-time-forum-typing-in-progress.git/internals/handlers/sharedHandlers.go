// Package handlers - shared helper utilities used by multiple handlers.
/*
Routes owned: none directly; used by handlers routed in routes.go.
Auth: helpers assume caller enforces auth; no auth checks here.
Side effects: DB reads only; no writes or HTTP responses.
Not responsible for: routing, JSON/HTML responses, or authorization decisions.
*/
package handlers

import (
	"database/sql"
	"realtimeforum/internals/utils"
)

// ============================================================================
// Section: Post helper lookups
// ============================================================================

// getPostTags returns the category tags for a given post.
// Caller supplies an open DB handle; expects valid postID; returns empty slice on query error.
func getPostTags(db *sql.DB, postID int) []string {
	rows, err := db.Query(`
		SELECT c.name FROM Categories c
		JOIN PostCategories pc ON c.category_id = pc.category_id
		WHERE pc.post_id = ?`, postID)
	if err != nil {
		return []string{}
	}
	defer rows.Close()

	var tags []string
	for rows.Next() {
		var tag string
		rows.Scan(&tag)
		tags = append(tags, tag)
	}
	return tags
}

// ============================================================================
// Section: Session helpers
// ============================================================================

// GetUsernameFromSession returns the username for a given session cookie.
// Delegates to utils without performing auth checks; callers must verify session validity.
func GetUsernameFromSession(cookieValue string) string {
	return utils.GetUsernameFromSession(cookieValue)
}
