package handlers

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/words-api/words/internal/models"
	"github.com/words-api/words/internal/services"
)

// ReviewHandler handles HTTP requests for review operations
type ReviewHandler struct {
	service *services.ReviewService
}

// NewReviewHandler creates a new review handler
func NewReviewHandler(db *sql.DB) *ReviewHandler {
	return &ReviewHandler{
		service: services.NewReviewService(db),
	}
}

// GetDueWords handles GET /api/users/:username/review
func (h *ReviewHandler) GetDueWords(c *gin.Context) {
	username := c.Param("username")

	dueWords, err := h.service.GetDueWords(username)
	if err != nil {
		if err.Error() == "user not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "user not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to retrieve due words",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"words": dueWords,
		"count": len(dueWords),
	})
}

// SubmitReview handles POST /api/users/:username/review/:word
func (h *ReviewHandler) SubmitReview(c *gin.Context) {
	username := c.Param("username")
	word := c.Param("word")

	var request models.ReviewRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "quality rating (0-5) is required",
		})
		return
	}

	updatedWord, err := h.service.SubmitReview(username, word, request.Quality)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "user not found" || err.Error() == "word not found" {
			statusCode = http.StatusNotFound
		} else if err.Error() == "quality must be between 0 and 5" {
			statusCode = http.StatusBadRequest
		}

		c.JSON(statusCode, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, updatedWord)
}

// GetReviewHistory handles GET /api/users/:username/review/:word/history
func (h *ReviewHandler) GetReviewHistory(c *gin.Context) {
	username := c.Param("username")
	word := c.Param("word")

	history, err := h.service.GetReviewHistory(username, word)
	if err != nil {
		if err.Error() == "user not found" || err.Error() == "word not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to retrieve review history",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"history": history,
		"count":   len(history),
	})
}
