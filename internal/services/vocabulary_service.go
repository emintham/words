package services

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/words-api/words/internal/models"
)

// VocabularyService handles business logic for vocabulary tracking
type VocabularyService struct {
	db          *sql.DB
	userService *UserService
	wordService *WordService
}

// NewVocabularyService creates a new vocabulary service
func NewVocabularyService(db *sql.DB) *VocabularyService {
	return &VocabularyService{
		db:          db,
		userService: NewUserService(db),
		wordService: NewWordService(db),
	}
}

// AddWord adds a word to a user's study list
func (s *VocabularyService) AddWord(username, wordStr string) (*models.UserWord, error) {
	// Get user
	user, err := s.userService.GetUser(username)
	if err != nil {
		return nil, err
	}

	// Normalize word
	wordStr = strings.ToLower(strings.TrimSpace(wordStr))
	if wordStr == "" {
		return nil, fmt.Errorf("word cannot be empty")
	}

	// Ensure word exists in the words table (fetch if needed)
	word, err := s.wordService.GetWord(wordStr)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch word: %w", err)
	}

	// Check if word is already in user's list
	var existingID int64
	err = s.db.QueryRow(`
		SELECT id FROM user_words WHERE user_id = ? AND word_id = ?
	`, user.ID, word.ID).Scan(&existingID)

	if err == nil {
		// Word already exists, return existing record
		return s.GetUserWord(user.ID, word.ID)
	}

	// Add word to user's study list
	nextReview := time.Now().Add(1 * time.Hour) // First review in 1 hour
	result, err := s.db.Exec(`
		INSERT INTO user_words (user_id, word_id, added_at, status, next_review_date, ease_factor, interval_days)
		VALUES (?, ?, ?, 'learning', ?, 2.5, 1)
	`, user.ID, word.ID, time.Now(), nextReview)

	if err != nil {
		return nil, fmt.Errorf("failed to add word: %w", err)
	}

	userWordID, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get user word ID: %w", err)
	}

	// Retrieve and return the created record
	return s.GetUserWordByID(userWordID)
}

// GetUserWords retrieves all words for a user, optionally filtered by status
func (s *VocabularyService) GetUserWords(username, status string) ([]models.UserWord, error) {
	// Get user
	user, err := s.userService.GetUser(username)
	if err != nil {
		return nil, err
	}

	query := `
		SELECT uw.id, uw.user_id, uw.word_id, w.word, uw.added_at, uw.status,
		       uw.next_review_date, uw.ease_factor, uw.interval_days
		FROM user_words uw
		JOIN words w ON uw.word_id = w.id
		WHERE uw.user_id = ?
	`
	args := []interface{}{user.ID}

	// Filter by status if provided
	if status != "" {
		query += " AND uw.status = ?"
		args = append(args, status)
	}

	query += " ORDER BY uw.added_at DESC"

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get user words: %w", err)
	}
	defer rows.Close()

	var userWords []models.UserWord
	for rows.Next() {
		var uw models.UserWord
		err := rows.Scan(&uw.ID, &uw.UserID, &uw.WordID, &uw.Word, &uw.AddedAt,
			&uw.Status, &uw.NextReviewDate, &uw.EaseFactor, &uw.IntervalDays)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user word: %w", err)
		}
		userWords = append(userWords, uw)
	}

	return userWords, nil
}

// GetUserWord retrieves a specific user word
func (s *VocabularyService) GetUserWord(userID, wordID int64) (*models.UserWord, error) {
	uw := &models.UserWord{}
	err := s.db.QueryRow(`
		SELECT uw.id, uw.user_id, uw.word_id, w.word, uw.added_at, uw.status,
		       uw.next_review_date, uw.ease_factor, uw.interval_days
		FROM user_words uw
		JOIN words w ON uw.word_id = w.id
		WHERE uw.user_id = ? AND uw.word_id = ?
	`, userID, wordID).Scan(&uw.ID, &uw.UserID, &uw.WordID, &uw.Word, &uw.AddedAt,
		&uw.Status, &uw.NextReviewDate, &uw.EaseFactor, &uw.IntervalDays)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user word not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user word: %w", err)
	}

	return uw, nil
}

// GetUserWordByID retrieves a user word by its ID
func (s *VocabularyService) GetUserWordByID(id int64) (*models.UserWord, error) {
	uw := &models.UserWord{}
	err := s.db.QueryRow(`
		SELECT uw.id, uw.user_id, uw.word_id, w.word, uw.added_at, uw.status,
		       uw.next_review_date, uw.ease_factor, uw.interval_days
		FROM user_words uw
		JOIN words w ON uw.word_id = w.id
		WHERE uw.id = ?
	`, id).Scan(&uw.ID, &uw.UserID, &uw.WordID, &uw.Word, &uw.AddedAt,
		&uw.Status, &uw.NextReviewDate, &uw.EaseFactor, &uw.IntervalDays)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user word not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user word: %w", err)
	}

	return uw, nil
}
