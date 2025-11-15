package auth

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"sync"
	"time"

	"github.com/words-api/words/internal/models"
)

var (
	ErrSessionNotFound = errors.New("session not found")
	ErrSessionExpired  = errors.New("session expired")
)

// Session represents an active user session
type Session struct {
	Token     string
	User      *models.User
	CreatedAt time.Time
	ExpiresAt time.Time
}

// SessionStore manages active sessions
type SessionStore struct {
	sessions map[string]*Session
	mu       sync.RWMutex
}

// NewSessionStore creates a new session store
func NewSessionStore() *SessionStore {
	store := &SessionStore{
		sessions: make(map[string]*Session),
	}
	// Start cleanup goroutine to remove expired sessions
	go store.cleanupExpiredSessions()
	return store
}

// CreateSession creates a new session for the given user
func (s *SessionStore) CreateSession(user *models.User) (*Session, error) {
	token, err := generateSessionToken()
	if err != nil {
		return nil, err
	}

	session := &Session{
		Token:     token,
		User:      user,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(24 * time.Hour), // 24-hour sessions
	}

	s.mu.Lock()
	s.sessions[token] = session
	s.mu.Unlock()

	return session, nil
}

// GetSession retrieves a session by token
func (s *SessionStore) GetSession(token string) (*Session, error) {
	s.mu.RLock()
	session, exists := s.sessions[token]
	s.mu.RUnlock()

	if !exists {
		return nil, ErrSessionNotFound
	}

	if time.Now().After(session.ExpiresAt) {
		s.DeleteSession(token)
		return nil, ErrSessionExpired
	}

	return session, nil
}

// DeleteSession removes a session by token
func (s *SessionStore) DeleteSession(token string) {
	s.mu.Lock()
	delete(s.sessions, token)
	s.mu.Unlock()
}

// cleanupExpiredSessions periodically removes expired sessions
func (s *SessionStore) cleanupExpiredSessions() {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	for range ticker.C {
		s.mu.Lock()
		now := time.Now()
		for token, session := range s.sessions {
			if now.After(session.ExpiresAt) {
				delete(s.sessions, token)
			}
		}
		s.mu.Unlock()
	}
}

// generateSessionToken generates a cryptographically secure random token
func generateSessionToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
