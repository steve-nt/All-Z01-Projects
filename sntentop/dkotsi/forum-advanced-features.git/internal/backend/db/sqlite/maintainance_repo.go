package sqlite

import (
	"database/sql"
	"fmt"
	"forum-advanced-features/internal/backend/models"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	DB      *sql.DB
	Config  *models.Config
	logfile *os.File
}

func NewDatabase(conf *models.Config, logfile *os.File) *sql.DB {
	db := &Database{Config: conf, logfile: logfile}
	db.Initialize()
	return db.DB
}

func (db *Database) Initialize() {
	dbPath := "../../data/forum.db"
	if dbPath == "" {
		dbPath = "../../data/forum.db"
	}

	isNewDB := false
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		file, err := os.Create(dbPath)
		if err != nil {
			log.Fatalf("failed to create database file: %v", err)
		}
		file.Close()
		isNewDB = true
		log.Println("Database file created successfully.")
	}

	newdb, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal(err)
	}
	newdb.Exec("PRAGMA foreign_keys=ON;")
	db.DB = newdb

	if isNewDB {
		migrationPath := "../../assets/migrations/posts.sql"
		if err := db.runMigration(migrationPath); err != nil {
			log.Fatalf("migration failed: %v", err)
		}
	}

	if err := db.CleanUpDatabase(); err != nil {
		log.Println(err)
	}

	log.Println("Database initialized successfully.")
}

func (db *Database) CleanUpDatabase() error {

	DbCleanUpRate := time.Duration(time.Hour) * time.Duration(db.Config.Durations.DbCleanUpRate)
	ticker := time.NewTicker(DbCleanUpRate)

	go func() error {
		for range ticker.C {
			db.Maintain()
		}
		return nil
	}()

	return nil
}

func (db *Database) runMigration(path string) error {
	content, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		return fmt.Errorf("failed to read migration file: %w", err)
	}

	// Διαχωρισμός των queries με ';'
	queries := strings.Split(string(content), ";")
	for _, query := range queries {
		query = strings.TrimSpace(query)
		if query == "" {
			continue
		}
		if _, err := db.DB.Exec(query); err != nil {
			return fmt.Errorf("failed to execute query: %w\nQuery: %s", err, query)
		}
	}
	log.Println("Migration executed successfully.")
	return nil
}

func (db *Database) Maintain() error {
	query := "DELETE FROM sessions WHERE expiration < ?;"

	_, err := db.DB.Exec(query, time.Now())
	if err != nil {
		log.Println(err)
		return err
	}
	db.CopyToBackUp()
	if err := db.logfile.Truncate(0); err != nil {
		log.Println(err)
	}
	return nil
}
func (db *Database) CopyToBackUp() {
	command := exec.Command("cp", "../../forum.db", "../../back-up/backup.db")
	err := command.Run()
	if err != nil {
		log.Println(err)
	}
}
