package handlers

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/words-api/words/internal/auth"
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

// GetDueWords handles GET /api/review (authenticated endpoint)
func (h *ReviewHandler) GetDueWords(c *gin.Context) {
	user, exists := auth.GetAuthenticatedUser(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "not authenticated",
		})
		return
	}

	dueWords, err := h.service.GetDueWords(user.Username)
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

// SubmitReview handles POST /api/review/:word (authenticated endpoint)
func (h *ReviewHandler) SubmitReview(c *gin.Context) {
	user, exists := auth.GetAuthenticatedUser(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "not authenticated",
		})
		return
	}

	word := c.Param("word")

	var request models.ReviewRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "quality rating (0-5) is required",
		})
		return
	}

	updatedWord, err := h.service.SubmitReview(user.Username, word, request.Quality)
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

// GetReviewHistory handles GET /api/review/:word/history (authenticated endpoint)
func (h *ReviewHandler) GetReviewHistory(c *gin.Context) {
	user, exists := auth.GetAuthenticatedUser(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "not authenticated",
		})
		return
	}

	word := c.Param("word")

	history, err := h.service.GetReviewHistory(user.Username, word)
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
