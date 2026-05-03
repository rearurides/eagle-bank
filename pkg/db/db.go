package db

import (
	"database/sql"
	"log"

	_ "modernc.org/sqlite"
)

func NewSQLiteDB(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	// SQLite performance settings
	_, err = db.Exec(`
		PRAGMA journal_mode=WAL;
		PRAGMA foreign_keys=ON;
		PRAGMA busy_timeout=5000;
	`)
	if err != nil {
		return nil, err
	}

	log.Println("Connected to SQLite database")
	return db, nil
}
