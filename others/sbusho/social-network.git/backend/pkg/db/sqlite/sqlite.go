package sqlite

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file" // File source driver for migrations
	_ "github.com/mattn/go-sqlite3"                      // SQLite driver - the underscore means we import for side effects only
)

// DB is a global variable to hold our database connection
// This allows other parts of the application to access the database
var DB *sql.DB

// InitDB initializes the SQLite database connection and runs migrations
//
// What this function does:
// 1. Creates the database file if it doesn't exist
// 2. Opens a connection to the SQLite database
// 3. Applies all pending migrations (creates/updates tables)
// 4. Stores the connection in the global DB variable
//
// Parameters:
//   - dbPath: The file path where the SQLite database will be stored
//     Example: "data/social_network.db"
//
// Returns:
//   - error: Any error that occurred during initialization
func InitDB(dbPath string) error {
	// Step 1: Ensure the directory for the database file exists
	// This prevents errors if the directory doesn't exist
	dbDir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		return fmt.Errorf("failed to create database directory: %w", err)
	}

	// Step 2: Open connection to SQLite database
	// The connection string format for SQLite is: "file:path?query_params"
	// _foreign_keys=1 enables foreign key constraints (important for data integrity)
	// _journal_mode=WAL enables Write-Ahead Logging (better performance for concurrent access)
	db, err := sql.Open("sqlite3", fmt.Sprintf("%s?_foreign_keys=1&_journal_mode=WAL", dbPath))
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	// Step 3: Test the connection by pinging the database
	// This ensures the database file is accessible and working
	if err := db.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	// Step 4: Store the connection in the global variable
	DB = db

	// Step 5: Run migrations to set up the database schema
	// Migrations will create all the tables defined in the migration files
	if err := runMigrations(db); err != nil {
		// If migrations fail, close the database connection
		db.Close()
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	log.Printf("Database initialized successfully at: %s", dbPath)
	return nil
}

// runMigrations applies all database migrations using golang-migrate
//
// What are migrations?
// Migrations are version-controlled database schema changes. They allow you to:
// - Track changes to your database structure over time
// - Apply changes in a controlled, repeatable way
// - Roll back changes if needed
// - Work with a team without database conflicts
//
// How golang-migrate works:
// 1. It looks for migration files in a directory (our migrations folder)
// 2. Each migration has an "up" file (applies changes) and "down" file (rolls back changes)
// 3. It tracks which migrations have been applied in a special table
// 4. It applies only the migrations that haven't been run yet
//
// Parameters:
//   - db: The database connection to apply migrations to
//
// Returns:
//   - error: Any error that occurred during migration
func runMigrations(db *sql.DB) error {
	// Step 1: Create a migration driver instance for SQLite
	// This tells golang-migrate how to interact with our SQLite database
	driver, err := sqlite3.WithInstance(db, &sqlite3.Config{})
	if err != nil {
		return fmt.Errorf("failed to create migration driver: %w", err)
	}

	// Step 2: Get the absolute path to the migrations directory
	// We need the absolute path because file:// URLs require it
	// Try relative path first (when running from backend/ directory)
	migrationsPath, err := filepath.Abs("pkg/db/migrations/sqlite")
	if err != nil {
		return fmt.Errorf("failed to get migrations path: %w", err)
	}

	// Check if the path exists, if not try with "backend/" prefix (when running from project root)
	if _, err := os.Stat(migrationsPath); os.IsNotExist(err) {
		// Try with backend/ prefix (for running from project root)
		altPath, err := filepath.Abs("backend/pkg/db/migrations/sqlite")
		if err == nil {
			if _, err := os.Stat(altPath); err == nil {
				migrationsPath = altPath
			}
		}
	}

	// Step 3: Create the migration source
	// The file:// protocol tells golang-migrate to read migration files from the filesystem
	// Format: file:///absolute/path/to/migrations
	migrationsURL := fmt.Sprintf("file:///%s", filepath.ToSlash(migrationsPath))

	// Step 4: Create a new migrate instance
	// This combines the database driver and the migration source
	m, err := migrate.NewWithDatabaseInstance(
		migrationsURL, // Where to find migration files
		"sqlite3",     // Database type
		driver,        // Database driver instance
	)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}

	// Step 5: Apply all pending migrations
	// Up() applies all migrations that haven't been run yet
	// If all migrations are already applied, it returns migrate.ErrNoChange
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	// Step 6: Get migration version info (for logging/debugging)
	version, dirty, err := m.Version()
	if err != nil && err != migrate.ErrNilVersion {
		return fmt.Errorf("failed to get migration version: %w", err)
	}

	if err == migrate.ErrNilVersion {
		log.Println("No migrations have been applied yet")
	} else {
		if dirty {
			log.Printf("Warning: Database is in a dirty state at version %d", version)
		} else {
			log.Printf("Migrations applied successfully. Current version: %d", version)
		}
	}

	return nil
}

// Insert executes a simple INSERT statement using positional placeholders based on args.
// Example: Insert(db, "Sessions", "(user_id, cookie_value)", 1, "abc")
func Insert(db *sql.DB, table, columns string, args ...interface{}) (sql.Result, error) {
	placeholders := make([]string, len(args))
	for i := range placeholders {
		placeholders[i] = "?"
	}
	query := fmt.Sprintf("INSERT INTO %s %s VALUES (%s)", table, columns, strings.Join(placeholders, ", "))
	return db.Exec(query, args...)
}

// GetDB returns the global database connection.
// Panics if InitDB has not been called.
func GetDB() *sql.DB {
	if DB == nil {
		panic("database not initialized. Call InitDB() first")
	}
	return DB
}

// CloseDB closes the database connection
// This should be called when the application is shutting down
// to properly clean up resources
func CloseDB() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}

// Ping tests the database connection
// Useful for health checks or verifying the connection is still alive
//
// Returns:
//   - error: Error if the connection is not working
func Ping() error {
	if DB == nil {
		return fmt.Errorf("database not initialized")
	}
	return DB.Ping()
}
