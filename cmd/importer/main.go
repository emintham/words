package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/words-api/words/internal/database"
	"github.com/words-api/words/internal/models"
)

// WordsetEntry represents the Wordset JSON structure
type WordsetEntry struct {
	Word      string            `json:"word"`
	WordsetID string            `json:"wordset_id"`
	Meanings  []WordsetMeaning  `json:"meanings"`
}

type WordsetMeaning struct {
	ID         string   `json:"id"`
	Definition string   `json:"def"`
	Example    string   `json:"example,omitempty"`
	SpeechPart string   `json:"speech_part"`
	Synonyms   []string `json:"synonyms,omitempty"`
}

// WordIndex tracks words for deduplication
type WordIndex struct {
	words map[string]*models.Word
}

func NewWordIndex() *WordIndex {
	return &WordIndex{
		words: make(map[string]*models.Word),
	}
}

// AddOrMerge adds a word or merges meanings if it already exists
func (idx *WordIndex) AddOrMerge(word *models.Word) {
	key := strings.ToLower(word.Word)

	if existing, exists := idx.words[key]; exists {
		// Merge meanings - add new meanings that don't duplicate
		existing.Meanings = append(existing.Meanings, word.Meanings...)
		// Update phonetic if the new one has more info
		if word.Phonetic != "" && len(word.Phonetic) > len(existing.Phonetic) {
			existing.Phonetic = word.Phonetic
		}
		// Merge source URLs
		existing.SourceUrls = append(existing.SourceUrls, word.SourceUrls...)
		// Merge phonetics
		existing.Phonetics = append(existing.Phonetics, word.Phonetics...)
	} else {
		idx.words[key] = word
	}
}

func (idx *WordIndex) GetAll() []*models.Word {
	result := make([]*models.Word, 0, len(idx.words))
	for _, word := range idx.words {
		result = append(result, word)
	}
	return result
}

func convertWordsetToModel(entry WordsetEntry) *models.Word {
	word := &models.Word{
		Word:       strings.ToLower(entry.Word),
		Phonetic:   "", // Wordset doesn't include phonetics
		SourceUrls: []string{"https://github.com/wordset/wordset-dictionary"},
		Meanings:   []models.Meaning{},
	}

	// Group meanings by part of speech
	meaningsByPOS := make(map[string]*models.Meaning)

	for _, wm := range entry.Meanings {
		pos := wm.SpeechPart

		// Get or create meaning for this POS
		meaning, exists := meaningsByPOS[pos]
		if !exists {
			meaning = &models.Meaning{
				PartOfSpeech: pos,
				Definitions:  []models.Definition{},
				Synonyms:     []string{},
			}
			meaningsByPOS[pos] = meaning
		}

		// Add definition
		def := models.Definition{
			Definition: wm.Definition,
			Example:    wm.Example,
			Synonyms:   wm.Synonyms,
		}
		meaning.Definitions = append(meaning.Definitions, def)
	}

	// Convert map to slice
	for _, meaning := range meaningsByPOS {
		word.Meanings = append(word.Meanings, *meaning)
	}

	return word
}

func loadWordsetFile(filepath string, index *WordIndex) (int, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return 0, err
	}

	var entries map[string]WordsetEntry
	if err := json.Unmarshal(data, &entries); err != nil {
		return 0, err
	}

	count := 0
	for _, entry := range entries {
		word := convertWordsetToModel(entry)
		index.AddOrMerge(word)
		count++
	}

	return count, nil
}

func importToDatabase(db *sql.DB, words []*models.Word, batchSize int) error {
	total := len(words)
	imported := 0
	errors := 0
	startTime := time.Now()

	for i := 0; i < total; i += batchSize {
		end := i + batchSize
		if end > total {
			end = total
		}

		batch := words[i:end]

		for _, word := range batch {
			// Set timestamps
			now := time.Now()
			word.CreatedAt = now
			word.UpdatedAt = now

			// Use the same saveToDB logic from word_service
			if err := saveWordToDB(db, word); err != nil {
				log.Printf("Error importing '%s': %v", word.Word, err)
				errors++
				continue
			}
			imported++
		}

		// Progress update
		elapsed := time.Since(startTime)
		percentage := float64(imported) / float64(total) * 100
		rate := float64(imported) / elapsed.Seconds()
		remaining := time.Duration(float64(total-imported)/rate) * time.Second

		fmt.Printf("\rProgress: %d/%d (%.1f%%) | Rate: %.0f words/sec | ETA: %s | Errors: %d",
			imported, total, percentage, rate, remaining.Round(time.Second), errors)
	}

	fmt.Println()
	return nil
}

func saveWordToDB(db *sql.DB, word *models.Word) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Check if word already exists
	var existingID int64
	err = tx.QueryRow("SELECT id FROM words WHERE word = ?", word.Word).Scan(&existingID)
	if err == nil {
		// Word exists, skip
		return nil
	}

	// Insert word
	result, err := tx.Exec(`
		INSERT INTO words (word, phonetic, created_at, updated_at)
		VALUES (?, ?, ?, ?)
	`, word.Word, word.Phonetic, word.CreatedAt, word.UpdatedAt)
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

		// Insert meaning-level synonyms
		for _, syn := range m.Synonyms {
			_, err := tx.Exec(`
				INSERT INTO synonyms (meaning_id, synonym) VALUES (?, ?)
			`, meaningID, syn)
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

			// Insert definition-level synonyms
			for _, syn := range d.Synonyms {
				_, err := tx.Exec(`
					INSERT INTO synonyms (definition_id, synonym) VALUES (?, ?)
				`, definitionID, syn)
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

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: import <wordset-data-dir>")
		fmt.Println("Example: import datasets/wordset-dictionary-master/data")
		os.Exit(1)
	}

	dataDir := os.Args[1]

	// Initialize database
	db, err := database.InitDB("words.db")
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	fmt.Println("📚 Starting dictionary import...")
	fmt.Printf("📁 Source: %s\n", dataDir)

	// Phase 1: Load all files into index (with deduplication)
	fmt.Println("\n🔄 Phase 1: Loading and deduplicating...")
	index := NewWordIndex()

	files, err := filepath.Glob(filepath.Join(dataDir, "*.json"))
	if err != nil {
		log.Fatalf("Failed to list files: %v", err)
	}

	totalLoaded := 0
	for i, file := range files {
		count, err := loadWordsetFile(file, index)
		if err != nil {
			log.Printf("Error loading %s: %v", file, err)
			continue
		}
		totalLoaded += count
		fmt.Printf("  [%d/%d] %s: %d entries\n", i+1, len(files), filepath.Base(file), count)
	}

	words := index.GetAll()
	fmt.Printf("\n✅ Loaded %d entries, deduplicated to %d unique words\n", totalLoaded, len(words))

	// Phase 2: Import to database
	fmt.Println("\n🔄 Phase 2: Importing to database...")
	startTime := time.Now()

	if err := importToDatabase(db, words, 100); err != nil {
		log.Fatalf("Import failed: %v", err)
	}

	elapsed := time.Since(startTime)
	fmt.Printf("\n✅ Import complete in %s\n", elapsed.Round(time.Second))
	fmt.Printf("📊 Final stats: %d words in database\n", len(words))
}
