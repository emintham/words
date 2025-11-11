package services

import (
	"database/sql"
	"fmt"
	"math"
	"time"

	"github.com/words-api/words/internal/models"
)

// ReviewService handles business logic for spaced repetition reviews
type ReviewService struct {
	db               *sql.DB
	userService      *UserService
	vocabularyService *VocabularyService
}

// NewReviewService creates a new review service
func NewReviewService(db *sql.DB) *ReviewService {
	return &ReviewService{
		db:               db,
		userService:      NewUserService(db),
		vocabularyService: NewVocabularyService(db),
	}
}

// GetDueWords retrieves words that are due for review
func (s *ReviewService) GetDueWords(username string) ([]models.UserWord, error) {
	// Get user
	user, err := s.userService.GetUser(username)
	if err != nil {
		return nil, err
	}

	// Query words due for review (next_review_date <= now)
	rows, err := s.db.Query(`
		SELECT uw.id, uw.user_id, uw.word_id, w.word, uw.added_at, uw.status,
		       uw.next_review_date, uw.ease_factor, uw.interval_days
		FROM user_words uw
		JOIN words w ON uw.word_id = w.id
		WHERE uw.user_id = ? AND uw.next_review_date <= datetime('now')
		ORDER BY uw.next_review_date ASC
	`, user.ID)

	if err != nil {
		return nil, fmt.Errorf("failed to get due words: %w", err)
	}
	defer rows.Close()

	var dueWords []models.UserWord
	for rows.Next() {
		var uw models.UserWord
		err := rows.Scan(&uw.ID, &uw.UserID, &uw.WordID, &uw.Word, &uw.AddedAt,
			&uw.Status, &uw.NextReviewDate, &uw.EaseFactor, &uw.IntervalDays)
		if err != nil {
			return nil, fmt.Errorf("failed to scan due word: %w", err)
		}
		dueWords = append(dueWords, uw)
	}

	return dueWords, nil
}

// SubmitReview processes a review and updates the user's progress using SM-2 algorithm
func (s *ReviewService) SubmitReview(username, wordStr string, quality int) (*models.UserWord, error) {
	// Validate quality rating
	if quality < 0 || quality > 5 {
		return nil, fmt.Errorf("quality must be between 0 and 5")
	}

	// Get user
	user, err := s.userService.GetUser(username)
	if err != nil {
		return nil, err
	}

	// Get word ID
	word, err := s.vocabularyService.wordService.GetWord(wordStr)
	if err != nil {
		return nil, fmt.Errorf("word not found: %w", err)
	}

	// Get current user word data
	userWord, err := s.vocabularyService.GetUserWord(user.ID, word.ID)
	if err != nil {
		return nil, fmt.Errorf("word not in user's study list: %w", err)
	}

	// Calculate new values using SM-2 algorithm
	newEaseFactor, newInterval, newStatus := s.calculateSM2(
		userWord.EaseFactor,
		userWord.IntervalDays,
		quality,
		userWord.Status,
	)

	// Calculate next review date
	nextReview := time.Now().Add(time.Duration(newInterval) * 24 * time.Hour)

	// Start transaction
	tx, err := s.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Update user_words table
	_, err = tx.Exec(`
		UPDATE user_words
		SET ease_factor = ?, interval_days = ?, next_review_date = ?, status = ?
		WHERE id = ?
	`, newEaseFactor, newInterval, nextReview, newStatus, userWord.ID)

	if err != nil {
		return nil, fmt.Errorf("failed to update user word: %w", err)
	}

	// Insert into review_history
	_, err = tx.Exec(`
		INSERT INTO review_history (user_id, word_id, reviewed_at, quality, interval_days, ease_factor)
		VALUES (?, ?, ?, ?, ?, ?)
	`, user.ID, word.ID, time.Now(), quality, newInterval, newEaseFactor)

	if err != nil {
		return nil, fmt.Errorf("failed to insert review history: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	// Return updated user word
	return s.vocabularyService.GetUserWord(user.ID, word.ID)
}

// calculateSM2 implements the SM-2 (SuperMemo 2) algorithm
// Returns: (newEaseFactor, newIntervalDays, newStatus)
func (s *ReviewService) calculateSM2(currentEF float64, currentInterval int, quality int, currentStatus string) (float64, int, string) {
	var newEF float64
	var newInterval int
	var newStatus string

	// Calculate new ease factor
	// EF' = EF + (0.1 - (5 - q) * (0.08 + (5 - q) * 0.02))
	newEF = currentEF + (0.1 - float64(5-quality)*(0.08+float64(5-quality)*0.02))

	// Minimum ease factor is 1.3
	if newEF < 1.3 {
		newEF = 1.3
	}

	// Calculate new interval based on quality
	if quality < 3 {
		// Failed recall - reset to 1 day
		newInterval = 1
		newStatus = "learning"
	} else {
		// Successful recall
		if currentInterval == 1 {
			// Second review after successful first
			newInterval = 6
		} else {
			// Subsequent reviews
			newInterval = int(math.Round(float64(currentInterval) * newEF))
		}

		// Update status based on interval
		if newInterval < 21 {
			newStatus = "reviewing"
		} else {
			newStatus = "mastered"
		}
	}

	// Ensure minimum interval is 1 day
	if newInterval < 1 {
		newInterval = 1
	}

	return newEF, newInterval, newStatus
}

// GetReviewHistory retrieves review history for a user's word
func (s *ReviewService) GetReviewHistory(username, wordStr string) ([]models.ReviewHistory, error) {
	// Get user
	user, err := s.userService.GetUser(username)
	if err != nil {
		return nil, err
	}

	// Get word ID
	word, err := s.vocabularyService.wordService.GetWord(wordStr)
	if err != nil {
		return nil, fmt.Errorf("word not found: %w", err)
	}

	// Query review history
	rows, err := s.db.Query(`
		SELECT rh.id, rh.user_id, rh.word_id, w.word, rh.reviewed_at, rh.quality,
		       rh.interval_days, rh.ease_factor
		FROM review_history rh
		JOIN words w ON rh.word_id = w.id
		WHERE rh.user_id = ? AND rh.word_id = ?
		ORDER BY rh.reviewed_at DESC
	`, user.ID, word.ID)

	if err != nil {
		return nil, fmt.Errorf("failed to get review history: %w", err)
	}
	defer rows.Close()

	var history []models.ReviewHistory
	for rows.Next() {
		var rh models.ReviewHistory
		err := rows.Scan(&rh.ID, &rh.UserID, &rh.WordID, &rh.Word, &rh.ReviewedAt,
			&rh.Quality, &rh.IntervalDays, &rh.EaseFactor)
		if err != nil {
			return nil, fmt.Errorf("failed to scan review history: %w", err)
		}
		history = append(history, rh)
	}

	return history, nil
}
