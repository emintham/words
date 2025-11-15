package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/words-api/words/internal/models"
)

const (
	SessionCookieName = "words_session"
	UserContextKey    = "user"
)

// AuthMiddleware creates a middleware that requires authentication
func AuthMiddleware(store *SessionStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get session token from cookie
		token, err := c.Cookie(SessionCookieName)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "authentication required",
			})
			c.Abort()
			return
		}

		// Validate session
		session, err := store.GetSession(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid or expired session",
			})
			c.Abort()
			return
		}

		// Add user to context
		c.Set(UserContextKey, session.User)
		c.Next()
	}
}

// GetAuthenticatedUser retrieves the authenticated user from the context
func GetAuthenticatedUser(c *gin.Context) (*models.User, bool) {
	user, exists := c.Get(UserContextKey)
	if !exists {
		return nil, false
	}
	return user.(*models.User), true
}
