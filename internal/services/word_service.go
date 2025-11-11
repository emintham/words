package services

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/words-api/words/internal/models"
	"github.com/words-api/words/pkg/dictionary"
)

// WordService handles business logic for word operations
type WordService struct {
	db       *sql.DB
	dictAPI  *dictionary.Client
}

// NewWordService creates a new word service
func NewWordService(db *sql.DB) *WordService {
	return &WordService{
		db:      db,
		dictAPI: dictionary.NewClient(),
	}
}

// GetWord retrieves a word from local DB first, falls back to API if not found
func (s *WordService) GetWord(word string) (*models.Word, error) {
	word = strings.ToLower(word)

	// Try local database first
	localWord, err := s.getFromDB(word)
	if err == nil {
		fmt.Printf("✓ Cache hit: '%s' (served from local DB)\n", word)
		return localWord, nil
	}

	// If not found locally, fetch from API
	if err == sql.ErrNoRows {
		fmt.Printf("⚡ Cache miss: '%s' (fetching from API)\n", word)
		apiWord, err := s.dictAPI.FetchWord(word)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch from API: %w", err)
		}

		// Save to database for future lookups
		if err := s.saveToDB(apiWord); err != nil {
			// Log error but still return the word
			fmt.Printf("Warning: failed to save word to DB: %v\n", err)
		}

		return apiWord, nil
	}

	return nil, fmt.Errorf("failed to retrieve word: %w", err)
}

// getFromDB retrieves a word from the local database
func (s *WordService) getFromDB(word string) (*models.Word, error) {
	w := &models.Word{}

	// Get word basic info
	err := s.db.QueryRow(`
		SELECT id, word, phonetic, created_at, updated_at
		FROM words WHERE word = ?
	`, word).Scan(&w.ID, &w.Word, &w.Phonetic, &w.CreatedAt, &w.UpdatedAt)

	if err != nil {
		return nil, err
	}

	// Get phonetics
	phoneticRows, err := s.db.Query(`
		SELECT id, text, audio FROM phonetics WHERE word_id = ?
	`, w.ID)
	if err != nil {
		return nil, err
	}
	defer phoneticRows.Close()

	for phoneticRows.Next() {
		var p models.Phonetic
		if err := phoneticRows.Scan(&p.ID, &p.Text, &p.Audio); err != nil {
			return nil, err
		}
		p.WordID = w.ID
		w.Phonetics = append(w.Phonetics, p)
	}

	// Get meanings and their definitions
	meaningRows, err := s.db.Query(`
		SELECT id, part_of_speech FROM meanings WHERE word_id = ?
	`, w.ID)
	if err != nil {
		return nil, err
	}
	defer meaningRows.Close()

	for meaningRows.Next() {
		var m models.Meaning
		if err := meaningRows.Scan(&m.ID, &m.PartOfSpeech); err != nil {
			return nil, err
		}
		m.WordID = w.ID

		// Get definitions for this meaning
		defRows, err := s.db.Query(`
			SELECT id, definition, example FROM definitions WHERE meaning_id = ?
		`, m.ID)
		if err != nil {
			return nil, err
		}

		for defRows.Next() {
			var d models.Definition
			if err := defRows.Scan(&d.ID, &d.Definition, &d.Example); err != nil {
				defRows.Close()
				return nil, err
			}
			d.MeaningID = m.ID

			// Get definition-level synonyms and antonyms
			d.Synonyms, _ = s.getSynonyms(0, d.ID)
			d.Antonyms, _ = s.getAntonyms(0, d.ID)

			m.Definitions = append(m.Definitions, d)
		}
		defRows.Close()

		// Get meaning-level synonyms and antonyms
		m.Synonyms, _ = s.getSynonyms(m.ID, 0)
		m.Antonyms, _ = s.getAntonyms(m.ID, 0)

		w.Meanings = append(w.Meanings, m)
	}

	// Get source URLs
	urlRows, err := s.db.Query(`
		SELECT url FROM source_urls WHERE word_id = ?
	`, w.ID)
	if err != nil {
		return nil, err
	}
	defer urlRows.Close()

	for urlRows.Next() {
		var url string
		if err := urlRows.Scan(&url); err != nil {
			return nil, err
		}
		w.SourceUrls = append(w.SourceUrls, url)
	}

	return w, nil
}

// saveToDB saves a word to the local database
func (s *WordService) saveToDB(word *models.Word) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Insert word
	result, err := tx.Exec(`
		INSERT INTO words (word, phonetic, created_at, updated_at)
		VALUES (?, ?, ?, ?)
	`, word.Word, word.Phonetic, time.Now(), time.Now())
	if err != nil {
		return err
	}

	wordID, err := result.LastInsertId()
	if err != nil {
		return err
	}

	// Insert phonetics
	for _, p := range word.Phonetics {
		_, err := tx.Exec(`
			INSERT INTO phonetics (word_id, text, audio) VALUES (?, ?, ?)
		`, wordID, p.Text, p.Audio)
		if err != nil {
			return err
		}
	}

	// Insert meanings and definitions
	for _, m := range word.Meanings {
		result, err := tx.Exec(`
			INSERT INTO meanings (word_id, part_of_speech) VALUES (?, ?)
		`, wordID, m.PartOfSpeech)
		if err != nil {
			return err
		}

		meaningID, err := result.LastInsertId()
		if err != nil {
			return err
		}

		// Insert meaning-level synonyms and antonyms
		for _, syn := range m.Synonyms {
			_, err := tx.Exec(`
				INSERT INTO synonyms (meaning_id, synonym) VALUES (?, ?)
			`, meaningID, syn)
			if err != nil {
				return err
			}
		}

		for _, ant := range m.Antonyms {
			_, err := tx.Exec(`
				INSERT INTO antonyms (meaning_id, antonym) VALUES (?, ?)
			`, meaningID, ant)
			if err != nil {
				return err
			}
		}

		// Insert definitions
		for _, d := range m.Definitions {
			result, err := tx.Exec(`
				INSERT INTO definitions (meaning_id, definition, example) VALUES (?, ?, ?)
			`, meaningID, d.Definition, d.Example)
			if err != nil {
				return err
			}

			definitionID, err := result.LastInsertId()
			if err != nil {
				return err
			}

			// Insert definition-level synonyms and antonyms
			for _, syn := range d.Synonyms {
				_, err := tx.Exec(`
					INSERT INTO synonyms (definition_id, synonym) VALUES (?, ?)
				`, definitionID, syn)
				if err != nil {
					return err
				}
			}

			for _, ant := range d.Antonyms {
				_, err := tx.Exec(`
					INSERT INTO antonyms (definition_id, antonym) VALUES (?, ?)
				`, definitionID, ant)
				if err != nil {
					return err
				}
			}
		}
	}

	// Insert source URLs
	for _, url := range word.SourceUrls {
		_, err := tx.Exec(`
			INSERT INTO source_urls (word_id, url) VALUES (?, ?)
		`, wordID, url)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

// Helper functions to get synonyms and antonyms
func (s *WordService) getSynonyms(meaningID, definitionID int64) ([]string, error) {
	query := `SELECT synonym FROM synonyms WHERE `
	var args []interface{}

	if meaningID > 0 {
		query += `meaning_id = ? AND definition_id IS NULL`
		args = append(args, meaningID)
	} else {
		query += `definition_id = ?`
		args = append(args, definitionID)
	}

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var synonyms []string
	for rows.Next() {
		var syn string
		if err := rows.Scan(&syn); err != nil {
			return nil, err
		}
		synonyms = append(synonyms, syn)
	}

	return synonyms, nil
}

func (s *WordService) getAntonyms(meaningID, definitionID int64) ([]string, error) {
	query := `SELECT antonym FROM antonyms WHERE `
	var args []interface{}

	if meaningID > 0 {
		query += `meaning_id = ? AND definition_id IS NULL`
		args = append(args, meaningID)
	} else {
		query += `definition_id = ?`
		args = append(args, definitionID)
	}

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var antonyms []string
	for rows.Next() {
		var ant string
		if err := rows.Scan(&ant); err != nil {
			return nil, err
		}
		antonyms = append(antonyms, ant)
	}

	return antonyms, nil
}
