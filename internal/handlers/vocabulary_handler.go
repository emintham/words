package handlers

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/words-api/words/internal/auth"
	"github.com/words-api/words/internal/services"
)

// VocabularyHandler handles HTTP requests for vocabulary operations
type VocabularyHandler struct {
	service *services.VocabularyService
}

// NewVocabularyHandler creates a new vocabulary handler
func NewVocabularyHandler(db *sql.DB) *VocabularyHandler {
	return &VocabularyHandler{
		service: services.NewVocabularyService(db),
	}
}

// AddWord handles POST /api/words/:word (authenticated endpoint)
func (h *VocabularyHandler) AddWord(c *gin.Context) {
	user, exists := auth.GetAuthenticatedUser(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "not authenticated",
		})
		return
	}

	word := c.Param("word")

	userWord, err := h.service.AddWord(user.Username, word)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "user not found" {
			statusCode = http.StatusNotFound
		} else if err.Error() == "word cannot be empty" {
			statusCode = http.StatusBadRequest
		}

		c.JSON(statusCode, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, userWord)
}

// GetUserWords handles GET /api/words (authenticated endpoint)
func (h *VocabularyHandler) GetUserWords(c *gin.Context) {
	user, exists := auth.GetAuthenticatedUser(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "not authenticated",
		})
		return
	}

	status := c.Query("status") // Optional filter: learning, reviewing, mastered

	userWords, err := h.service.GetUserWords(user.Username, status)
	if err != nil {
		if err.Error() == "user not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "user not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to retrieve user words",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"words": userWords,
		"count": len(userWords),
	})
}
