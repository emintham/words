package models

import "time"

// User represents a user account
type User struct {
	ID        int64     `json:"id" db:"id"`
	Username  string    `json:"username" db:"username"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// UserWord represents a word that a user is studying
type UserWord struct {
	ID             int64     `json:"id" db:"id"`
	UserID         int64     `json:"user_id" db:"user_id"`
	WordID         int64     `json:"word_id" db:"word_id"`
	Word           string    `json:"word,omitempty"`
	AddedAt        time.Time `json:"added_at" db:"added_at"`
	Status         string    `json:"status" db:"status"` // learning, reviewing, mastered
	NextReviewDate time.Time `json:"next_review_date" db:"next_review_date"`
	EaseFactor     float64   `json:"ease_factor" db:"ease_factor"`
	IntervalDays   int       `json:"interval_days" db:"interval_days"`
}

// ReviewHistory represents a single review session
type ReviewHistory struct {
	ID           int64     `json:"id" db:"id"`
	UserID       int64     `json:"user_id" db:"user_id"`
	WordID       int64     `json:"word_id" db:"word_id"`
	Word         string    `json:"word,omitempty"`
	ReviewedAt   time.Time `json:"reviewed_at" db:"reviewed_at"`
	Quality      int       `json:"quality" db:"quality"` // 0-5 rating
	IntervalDays int       `json:"interval_days" db:"interval_days"`
	EaseFactor   float64   `json:"ease_factor" db:"ease_factor"`
}

// ReviewRequest represents the JSON body for submitting a review
type ReviewRequest struct {
	Quality int `json:"quality" binding:"required,min=0,max=5"`
}

// UserStats represents learning statistics for a user
type UserStats struct {
	Username       string    `json:"username"`
	TotalWords     int       `json:"total_words"`
	DueToday       int       `json:"due_today"`
	Learning       int       `json:"learning"`
	Reviewing      int       `json:"reviewing"`
	Mastered       int       `json:"mastered"`
	TotalReviews   int       `json:"total_reviews"`
	CurrentStreak  int       `json:"current_streak"`
	LastReviewDate time.Time `json:"last_review_date,omitempty"`
}
