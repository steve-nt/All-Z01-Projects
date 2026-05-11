package db

import "strings"

// Category represents one row from the categories table.
type Category struct {
	ID   int
	Name string
}

// GetAllCategories returns all categories sorted alphabetically.
// Used for:
// - Create Post checkboxes
// - Posts filter chips
func GetAllCategories() ([]Category, error) {
	rows, err := DB.Query(`SELECT id, name FROM categories ORDER BY name ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []Category

	for rows.Next() {
		var c Category
		if err := rows.Scan(&c.ID, &c.Name); err != nil {
			return nil, err
		}
		out = append(out, c)
	}

	return out, rows.Err()
}

// SetPostCategoryIDs replaces the category links for a post.
func SetPostCategoryIDs(postID int, categoryIDs []int) error {
	// Remove existing relations
	if _, err := DB.Exec(`DELETE FROM post_categories WHERE post_id = ?`, postID); err != nil {
		return err
	}

	// Insert selected categories
	for _, cid := range categoryIDs {

		if cid <= 0 {
			continue
		}

		if _, err := DB.Exec(
			`INSERT OR IGNORE INTO post_categories (post_id, category_id) VALUES (?, ?)`,
			postID,
			cid,
		); err != nil {
			return err
		}
	}

	return nil
}

// CreateCategory inserts a new category.
// AUDIT: category management belongs to administrators.
func CreateCategory(name string) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil
	}

	_, err := DB.Exec(`INSERT INTO categories (name) VALUES (?)`, name)
	return err
}

// DeleteCategory removes a category and any post-category relations using it.
// AUDIT: category deletion is intentionally explicit to keep forum taxonomy manageable.
func DeleteCategory(categoryID int) error {
	if _, err := DB.Exec(`DELETE FROM post_categories WHERE category_id = ?`, categoryID); err != nil {
		return err
	}

	_, err := DB.Exec(`DELETE FROM categories WHERE id = ?`, categoryID)
	return err
}