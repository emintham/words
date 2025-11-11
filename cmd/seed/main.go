package main

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", "words.db")
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	// Test words to insert
	testWords := []struct {
		word       string
		phonetic   string
		definition string
		partOfSpeech string
		example    string
	}{
		{"ephemeral", "/ɪˈfɛm(ə)ɹəl/", "Lasting for a very short time", "adjective", "The ephemeral beauty of cherry blossoms"},
		{"serendipity", "/ˌsɛɹ.ənˈdɪp.ɪ.ti/", "An unsought, unintended occurrence with fortunate results", "noun", "It was pure serendipity that we met"},
		{"mordant", "/ˈmɔːdənt/", "Biting and caustic in thought, manner, or style", "adjective", "Her mordant wit could cut through any pretense"},
		{"ubiquitous", "/juːˈbɪk.wɪ.təs/", "Being present everywhere at once", "adjective", "Smartphones have become ubiquitous in modern society"},
		{"ephemera", "/ɪˈfɛm.ə.ɹə/", "Items intended to be useful or important for only a short time", "noun", "The museum collected vintage ephemera like postcards"},
	}

	for _, w := range testWords {
		tx, err := db.Begin()
		if err != nil {
			log.Printf("Failed to begin transaction: %v", err)
			continue
		}

		// Insert word
		result, err := tx.Exec(`
			INSERT INTO words (word, phonetic, created_at, updated_at)
			VALUES (?, ?, ?, ?)
		`, w.word, w.phonetic, time.Now(), time.Now())
		if err != nil {
			tx.Rollback()
			log.Printf("Failed to insert word %s: %v", w.word, err)
			continue
		}

		wordID, _ := result.LastInsertId()

		// Insert meaning
		result, err = tx.Exec(`
			INSERT INTO meanings (word_id, part_of_speech)
			VALUES (?, ?)
		`, wordID, w.partOfSpeech)
		if err != nil {
			tx.Rollback()
			log.Printf("Failed to insert meaning for %s: %v", w.word, err)
			continue
		}

		meaningID, _ := result.LastInsertId()

		// Insert definition
		_, err = tx.Exec(`
			INSERT INTO definitions (meaning_id, definition, example)
			VALUES (?, ?, ?)
		`, meaningID, w.definition, w.example)
		if err != nil {
			tx.Rollback()
			log.Printf("Failed to insert definition for %s: %v", w.word, err)
			continue
		}

		if err := tx.Commit(); err != nil {
			log.Printf("Failed to commit transaction for %s: %v", w.word, err)
			continue
		}

		log.Printf("✓ Inserted word: %s", w.word)
	}

	log.Println("Seed data inserted successfully!")
}
