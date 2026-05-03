package db

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"
	"sort"
)

// RunMigrations applies all migration files in the specified path to the database
func RunMigrations(db *sql.DB, migrationsPath string) error {
	files, err := filepath.Glob(filepath.Join(migrationsPath, "*.up.sql"))
	if err != nil {
		return err
	}

	sort.Strings(files)

	for _, f := range files {
		content, err := os.ReadFile(f)
		if err != nil {
			return err
		}
		if _, err := db.Exec(string(content)); err != nil {
			return err
		}
		log.Printf("Applied migration: %s", filepath.Base(f))
	}
	return nil
}
