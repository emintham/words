package handlers

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/words-api/words/internal/auth"
	"github.com/words-api/words/internal/services"
)

// AuthHandler handles authentication operations
type AuthHandler struct {
	userService  *services.UserService
	sessionStore *auth.SessionStore
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(db *sql.DB, sessionStore *auth.SessionStore) *AuthHandler {
	return &AuthHandler{
		userService:  services.NewUserService(db),
		sessionStore: sessionStore,
	}
}

// Login handles POST /api/auth/login
func (h *AuthHandler) Login(c *gin.Context) {
	var request struct {
		Username string `json:"username" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "username is required",
		})
		return
	}

	// Try to get existing user
	user, err := h.userService.GetUser(request.Username)
	if err != nil {
		// If user doesn't exist, create new user
		if err.Error() == "user not found" {
			user, err = h.userService.CreateUser(request.Username)
			if err != nil {
				statusCode := http.StatusInternalServerError
				if err.Error() == "username already exists" {
					statusCode = http.StatusConflict
				} else if err.Error() == "username must be between 3 and 20 characters" ||
					err.Error() == "username can only contain letters, numbers, and underscores" {
					statusCode = http.StatusBadRequest
				}

				c.JSON(statusCode, gin.H{
					"error": err.Error(),
				})
				return
			}
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to authenticate",
			})
			return
		}
	}

	// Create session
	session, err := h.sessionStore.CreateSession(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to create session",
		})
		return
	}

	// Set session cookie
	c.SetCookie(
		auth.SessionCookieName,
		session.Token,
		int(24*time.Hour.Seconds()), // maxAge in seconds
		"/",                          // path
		"",                           // domain (empty for same-site)
		false,                        // secure (set to true in production with HTTPS)
		true,                         // httpOnly
	)

	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}

// Logout handles POST /api/auth/logout
func (h *AuthHandler) Logout(c *gin.Context) {
	// Get session token from cookie
	token, err := c.Cookie(auth.SessionCookieName)
	if err == nil {
		// Delete session from store
		h.sessionStore.DeleteSession(token)
	}

	// Clear cookie
	c.SetCookie(
		auth.SessionCookieName,
		"",
		-1,    // maxAge -1 deletes the cookie
		"/",   // path
		"",    // domain
		false, // secure
		true,  // httpOnly
	)

	c.JSON(http.StatusOK, gin.H{
		"message": "logged out successfully",
	})
}

// GetCurrentUser handles GET /api/auth/me
func (h *AuthHandler) GetCurrentUser(c *gin.Context) {
	user, exists := auth.GetAuthenticatedUser(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "not authenticated",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}
