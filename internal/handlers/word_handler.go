package handlers

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/words-api/words/internal/services"
)

// WordHandler handles HTTP requests for word operations
type WordHandler struct {
	service *services.WordService
}

// NewWordHandler creates a new word handler
func NewWordHandler(db *sql.DB) *WordHandler {
	return &WordHandler{
		service: services.NewWordService(db),
	}
}

// GetWord handles GET /api/words/:word
func (h *WordHandler) GetWord(c *gin.Context) {
	word := c.Param("word")

	if word == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "word parameter is required",
		})
		return
	}

	result, err := h.service.GetWord(word)
	if err != nil {
		if err.Error() == "word not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "word not found",
				"word":  word,
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to retrieve word",
		})
		return
	}

	c.JSON(http.StatusOK, result)
}
