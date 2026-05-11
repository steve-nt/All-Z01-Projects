package db

import (
	"strings"

	"forum/models"
)

// CreateUser inserts a new user into the database.
//
// Email is normalized to lowercase and trimmed so that
// "User@Mail.com" and "user@mail.com" are treated the same.
// Password must already be bcrypt-hashed before calling this.
func CreateUser(email, username, hashedPassword string) error {
	// Normalize input for consistency and duplicate detection.
	email = strings.TrimSpace(strings.ToLower(email))
	username = strings.TrimSpace(username)

	role := nextAssignedRole()

	// Insert user record.
	_, err := DB.Exec(
		`INSERT INTO users (email, username, password, role) VALUES (?, ?, ?, ?)`,
		email, username, hashedPassword, role,
	)

	return err
}

// GetUserByID fetches a user using their numeric ID.
// Used for logged-in UI details such as welcome messages and role checks.
func GetUserByID(id int) (models.User, error) {
	var u models.User

	err := DB.QueryRow(
		`SELECT id, email, username, password,
		        IFNULL(provider, ''), IFNULL(provider_id, ''), role
		 FROM users
		 WHERE id = ?`,
		id,
	).Scan(&u.ID, &u.Email, &u.Username, &u.Password, &u.Provider, &u.ProviderID, &u.Role)

	return u, err
}

// GetUserByEmail fetches a user using their email.
// Returns sql.ErrNoRows if the email does not exist.
func GetUserByEmail(email string) (models.User, error) {
	// Normalize email to match how it's stored.
	email = strings.TrimSpace(strings.ToLower(email))

	var u models.User

	err := DB.QueryRow(
		`SELECT id, email, username, password,
		        IFNULL(provider, ''), IFNULL(provider_id, ''), role
		 FROM users
		 WHERE email = ?`,
		email,
	).Scan(&u.ID, &u.Email, &u.Username, &u.Password, &u.Provider, &u.ProviderID, &u.Role)

	return u, err
}

// GetUserByUsername fetches a user using their username.
// Used to detect duplicate usernames during registration.
func GetUserByUsername(username string) (models.User, error) {
	// Trim whitespace for consistent comparison.
	username = strings.TrimSpace(username)

	var u models.User

	err := DB.QueryRow(
		`SELECT id, email, username, password,
		        IFNULL(provider, ''), IFNULL(provider_id, ''), role
		 FROM users
		 WHERE username = ?`,
		username,
	).Scan(&u.ID, &u.Email, &u.Username, &u.Password, &u.Provider, &u.ProviderID, &u.Role)

	return u, err
}