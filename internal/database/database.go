package database

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

// InitDB initializes the SQLite database and creates tables
func InitDB(filepath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Test connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Create tables
	if err := createTables(db); err != nil {
		return nil, fmt.Errorf("failed to create tables: %w", err)
	}

	return db, nil
}

func createTables(db *sql.DB) error {
	schema := `
	CREATE TABLE IF NOT EXISTS words (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		word TEXT NOT NULL UNIQUE,
		phonetic TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE INDEX IF NOT EXISTS idx_words_word ON words(word);

	CREATE TABLE IF NOT EXISTS phonetics (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		word_id INTEGER NOT NULL,
		text TEXT NOT NULL,
		audio TEXT,
		FOREIGN KEY (word_id) REFERENCES words(id) ON DELETE CASCADE
	);

	CREATE TABLE IF NOT EXISTS meanings (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		word_id INTEGER NOT NULL,
		part_of_speech TEXT NOT NULL,
		FOREIGN KEY (word_id) REFERENCES words(id) ON DELETE CASCADE
	);

	CREATE TABLE IF NOT EXISTS definitions (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		meaning_id INTEGER NOT NULL,
		definition TEXT NOT NULL,
		example TEXT,
		FOREIGN KEY (meaning_id) REFERENCES meanings(id) ON DELETE CASCADE
	);

	CREATE TABLE IF NOT EXISTS synonyms (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		meaning_id INTEGER,
		definition_id INTEGER,
		synonym TEXT NOT NULL,
		CHECK ((meaning_id IS NOT NULL) OR (definition_id IS NOT NULL)),
		FOREIGN KEY (meaning_id) REFERENCES meanings(id) ON DELETE CASCADE,
		FOREIGN KEY (definition_id) REFERENCES definitions(id) ON DELETE CASCADE
	);

	CREATE TABLE IF NOT EXISTS antonyms (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		meaning_id INTEGER,
		definition_id INTEGER,
		antonym TEXT NOT NULL,
		CHECK ((meaning_id IS NOT NULL) OR (definition_id IS NOT NULL)),
		FOREIGN KEY (meaning_id) REFERENCES meanings(id) ON DELETE CASCADE,
		FOREIGN KEY (definition_id) REFERENCES definitions(id) ON DELETE CASCADE
	);

	CREATE TABLE IF NOT EXISTS source_urls (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		word_id INTEGER NOT NULL,
		url TEXT NOT NULL,
		FOREIGN KEY (word_id) REFERENCES words(id) ON DELETE CASCADE
	);
	`

	_, err := db.Exec(schema)
	return err
}
