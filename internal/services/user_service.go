package services

import (
	"database/sql"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/words-api/words/internal/models"
)

// UserService handles business logic for user operations
type UserService struct {
	db *sql.DB
}

// NewUserService creates a new user service
func NewUserService(db *sql.DB) *UserService {
	return &UserService{db: db}
}

// CreateUser creates a new user with the given username
func (s *UserService) CreateUser(username string) (*models.User, error) {
	// Validate username
	username = strings.TrimSpace(username)
	if len(username) < 3 || len(username) > 20 {
		return nil, fmt.Errorf("username must be between 3 and 20 characters")
	}

	// Only allow alphanumeric characters and underscores
	validUsername := regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
	if !validUsername.MatchString(username) {
		return nil, fmt.Errorf("username can only contain letters, numbers, and underscores")
	}

	// Insert user into database
	_, err := s.db.Exec(`
		INSERT INTO users (username, created_at) VALUES (?, ?)
	`, username, time.Now())

	if err != nil {
		// Check if username already exists
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return nil, fmt.Errorf("username already exists")
		}
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Return the created user
	return s.GetUser(username)
}

// GetUser retrieves a user by username
func (s *UserService) GetUser(username string) (*models.User, error) {
	user := &models.User{}
	err := s.db.QueryRow(`
		SELECT id, username, created_at FROM users WHERE username = ?
	`, username).Scan(&user.ID, &user.Username, &user.CreatedAt)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

// GetUserByID retrieves a user by ID
func (s *UserService) GetUserByID(userID int64) (*models.User, error) {
	user := &models.User{}
	err := s.db.QueryRow(`
		SELECT id, username, created_at FROM users WHERE id = ?
	`, userID).Scan(&user.ID, &user.Username, &user.CreatedAt)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

// GetUserStats retrieves learning statistics for a user
func (s *UserService) GetUserStats(username string) (*models.UserStats, error) {
	user, err := s.GetUser(username)
	if err != nil {
		return nil, err
	}

	stats := &models.UserStats{
		Username: username,
	}

	// Get total words
	err = s.db.QueryRow(`
		SELECT COUNT(*) FROM user_words WHERE user_id = ?
	`, user.ID).Scan(&stats.TotalWords)
	if err != nil {
		return nil, fmt.Errorf("failed to get total words: %w", err)
	}

	// Get words due today
	err = s.db.QueryRow(`
		SELECT COUNT(*) FROM user_words
		WHERE user_id = ? AND next_review_date <= datetime('now')
	`, user.ID).Scan(&stats.DueToday)
	if err != nil {
		return nil, fmt.Errorf("failed to get due words: %w", err)
	}

	// Get count by status
	err = s.db.QueryRow(`
		SELECT COUNT(*) FROM user_words WHERE user_id = ? AND status = 'learning'
	`, user.ID).Scan(&stats.Learning)
	if err != nil {
		return nil, fmt.Errorf("failed to get learning count: %w", err)
	}

	err = s.db.QueryRow(`
		SELECT COUNT(*) FROM user_words WHERE user_id = ? AND status = 'reviewing'
	`, user.ID).Scan(&stats.Reviewing)
	if err != nil {
		return nil, fmt.Errorf("failed to get reviewing count: %w", err)
	}

	err = s.db.QueryRow(`
		SELECT COUNT(*) FROM user_words WHERE user_id = ? AND status = 'mastered'
	`, user.ID).Scan(&stats.Mastered)
	if err != nil {
		return nil, fmt.Errorf("failed to get mastered count: %w", err)
	}

	// Get total reviews
	err = s.db.QueryRow(`
		SELECT COUNT(*) FROM review_history WHERE user_id = ?
	`, user.ID).Scan(&stats.TotalReviews)
	if err != nil {
		return nil, fmt.Errorf("failed to get total reviews: %w", err)
	}

	// Get last review date and calculate streak
	var lastReviewDateStr sql.NullString
	err = s.db.QueryRow(`
		SELECT MAX(reviewed_at) FROM review_history WHERE user_id = ?
	`, user.ID).Scan(&lastReviewDateStr)
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("failed to get last review date: %w", err)
	}

	if lastReviewDateStr.Valid {
		// Parse the datetime string (RFC3339 format with nanoseconds)
		lastReviewDate, err := time.Parse(time.RFC3339Nano, lastReviewDateStr.String)
		if err == nil {
			stats.LastReviewDate = lastReviewDate
			// Calculate streak, but don't fail if it errors
			streak := s.calculateStreak(user.ID)
			stats.CurrentStreak = streak
		}
	}

	return stats, nil
}

// calculateStreak calculates the current consecutive days streak for a user
func (s *UserService) calculateStreak(userID int64) int {
	rows, err := s.db.Query(`
		SELECT DISTINCT strftime('%Y-%m-%d', reviewed_at) as review_date
		FROM review_history
		WHERE user_id = ?
		ORDER BY review_date DESC
	`, userID)
	if err != nil {
		return 0
	}
	defer rows.Close()

	var dates []time.Time
	for rows.Next() {
		var dateStr string
		if err := rows.Scan(&dateStr); err != nil {
			return 0
		}
		date, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			return 0
		}
		dates = append(dates, date)
	}

	if len(dates) == 0 {
		return 0
	}

	// Check if the most recent review was today or yesterday
	now := time.Now().UTC().Truncate(24 * time.Hour)
	mostRecent := dates[0].UTC().Truncate(24 * time.Hour)

	daysSince := int(now.Sub(mostRecent).Hours() / 24)
	if daysSince > 1 {
		return 0 // Streak is broken
	}

	// Count consecutive days
	streak := 1
	for i := 1; i < len(dates); i++ {
		current := dates[i].UTC().Truncate(24 * time.Hour)
		previous := dates[i-1].UTC().Truncate(24 * time.Hour)

		diff := int(previous.Sub(current).Hours() / 24)
		if diff == 1 {
			streak++
		} else {
			break
		}
	}

	return streak
}
