package db

import (
	"database/sql"
	"log"
	"os"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

// DB is the global database connection used by the whole application.
// It is initialized once in Init() and reused everywhere.
var DB *sql.DB

// Init opens the SQLite database, enables constraints,
// then ensures tables, indexes and default data exist.
func Init() {
	var err error

	// DB_PATH allows Docker to store the database in a mounted directory.
	// If not set, we default to a local file in the data directory.
	dbPath := strings.TrimSpace(os.Getenv("DB_PATH"))
	if dbPath == "" {
		// AUDIT: default database location moved to ./data to avoid polluting project root.
		dbPath = "./data/forum.db"
	}

	// AUDIT: ensure database directory exists before opening SQLite file.
	if err := os.MkdirAll("./data", 0755); err != nil {
		log.Fatal("could not create data directory:", err)
	}

	// Open SQLite database file (creates it if missing)
	DB, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal(err)
	}

	// SQLite disables foreign keys by default.
	// This enables relational integrity checks
	// (required for proper FK behaviour).
	if _, err := DB.Exec(`PRAGMA foreign_keys = ON;`); err != nil {
		log.Println("WARNING: could not enable foreign_keys:", err)
	}

	// Recommended settings for better concurrency and performance.
	_, _ = DB.Exec(`PRAGMA busy_timeout = 5000;`)
	_, _ = DB.Exec(`PRAGMA journal_mode = WAL;`)
	_, _ = DB.Exec(`PRAGMA synchronous = NORMAL;`)

	// Verify connection is actually working
	if err := DB.Ping(); err != nil {
		log.Fatal(err)
	}

	// Ensure DB structure exists
	createTables()

	// Apply schema migrations for databases created before newer features existed.
	if err := migrateUsersOAuthColumns(); err != nil {
		log.Fatal(err)
	}
	if err := migrateUserRoles(); err != nil {
		log.Fatal(err)
	}
	if err := migratePostModerationColumns(); err != nil {
		log.Fatal(err)
	}

	// AUDIT: moderation optional adds moderator-request workflow.
	if err := migrateModeratorRequestsTable(); err != nil {
		log.Fatal(err)
	}

	// AUDIT: moderation optional adds moderator -> admin reporting.
	if err := migrateReportsTable(); err != nil {
		log.Fatal(err)
	}

	// AUDIT: notifications must support generic moderation-related messages.
	if err := migrateNotificationsForModeration(); err != nil {
		log.Fatal(err)
	}

	// Ensure an administrator always exists.
	// For an older database with users but no role assignments yet,
	// the oldest user is promoted to admin so moderation remains manageable.
	if err := ensureAdminExists(); err != nil {
		log.Fatal(err)
	}

	// Ensure additional constraints exist
	createIndexes()

	// Insert default categories if DB is empty
	seedCategories()
}