package handlers

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
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

// AddWord handles POST /api/users/:username/words/:word
func (h *VocabularyHandler) AddWord(c *gin.Context) {
	username := c.Param("username")
	word := c.Param("word")

	userWord, err := h.service.AddWord(username, word)
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

// GetUserWords handles GET /api/users/:username/words
func (h *VocabularyHandler) GetUserWords(c *gin.Context) {
	username := c.Param("username")
	status := c.Query("status") // Optional filter: learning, reviewing, mastered

	userWords, err := h.service.GetUserWords(username, status)
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
