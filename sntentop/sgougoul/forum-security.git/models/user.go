package models

// User represents one user row from the database.
type User struct {
	ID         int
	Email      string
	Username   string
	Password   string
	Provider   string
	ProviderID string
	Role       string
}