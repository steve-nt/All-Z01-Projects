package models

import (
	"database/sql"
	"errors"
	"forum/src/utils"
	"os"
	"path"
	"runtime"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

const migrations_dir = "migrations"

var db *sql.DB

type MigrationsEnabled map[string]bool

func InitDB(dbPath string) error {
	var err error
	db, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		err = errors.Join(utils.GetFunctionName(), err)
		return err
	}
	_, err = db.Exec("PRAGMA journal_mode=WAL;")
	if err != nil {
		err = errors.Join(utils.GetFunctionName(), err)
		return err
	}
	db.SetMaxOpenConns(1)
	err = createMigrationsTable(db)
	if err != nil {
		err = errors.Join(utils.GetFunctionName(), err)
		return err
	}
	err = runMigrations(db)
	if err != nil {
		err = errors.Join(utils.GetFunctionName(), err)
		return err
	}
	return nil
}

func createMigrationsTable(db *sql.DB) error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS "schema_migrations" (
		"version"	TEXT NOT NULL UNIQUE,
		"timestamp"	TEXT NOT NULL DEFAULT current_timestamp
	)`)
	if err != nil {
		err = errors.Join(utils.GetFunctionName(), err)
		return err
	}
	return nil
}

func selectMigrations(db *sql.DB) (MigrationsEnabled, error) {
	applied := MigrationsEnabled{}
	var err error
	rows, err := db.Query("SELECT version FROM schema_migrations")
	if err != nil {
		err = errors.Join(utils.GetFunctionName(), err)
		return applied, err
	}
	defer rows.Close()
	for rows.Next() {
		var version string
		rows.Scan(&version)
		applied[version] = true
	}
	return applied, err
}

func runMigrations(db *sql.DB) error {
	applied := MigrationsEnabled{}
	_, x, _, ok := runtime.Caller(0)
	if !ok {
		return ErrorFailedToGetCaller
	}
	migrations_dir := path.Join(path.Dir(x), "..", "..", migrations_dir)
	utils.LogDebug(x)
	migrations_found, err := os.ReadDir(migrations_dir)
	if err != nil {
		err = errors.Join(utils.GetFunctionName(), err)
		return err
	}
	applied, err = selectMigrations(db)
	if err != nil {
		err = errors.Join(utils.GetFunctionName(), err)
		return err
	}
	for _, file := range migrations_found {
		if !applied[file.Name()] {
			bytes, err := os.ReadFile(path.Join(migrations_dir, file.Name()))
			if err != nil {
				err = errors.Join(utils.GetFunctionName(), err)
				return err
			}
			if strings.HasSuffix(file.Name(), ".sql") {
				utils.LogInfo("Running migration file: "+file.Name())
				query := string(bytes)
				_, err = db.Exec(query)
				if err != nil {
					err = errors.Join(utils.GetFunctionName(), err)
					return err
				}
				_, err = db.Exec("INSERT INTO schema_migrations(version) VALUES (?)", file.Name())
				if err != nil {
					err = errors.Join(utils.GetFunctionName(), err)
					return err
				}
			}
		}
	}
	return nil
}
