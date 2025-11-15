package handlers

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/words-api/words/internal/auth"
	"github.com/words-api/words/internal/services"
)

// UserHandler handles HTTP requests for user operations
type UserHandler struct {
	service *services.UserService
}

// NewUserHandler creates a new user handler
func NewUserHandler(db *sql.DB) *UserHandler {
	return &UserHandler{
		service: services.NewUserService(db),
	}
}

// CreateUser handles POST /api/users
func (h *UserHandler) CreateUser(c *gin.Context) {
	var request struct {
		Username string `json:"username" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "username is required",
		})
		return
	}

	user, err := h.service.CreateUser(request.Username)
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

	c.JSON(http.StatusCreated, user)
}

// GetUser handles GET /api/user (authenticated endpoint)
func (h *UserHandler) GetUser(c *gin.Context) {
	user, exists := auth.GetAuthenticatedUser(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "not authenticated",
		})
		return
	}

	c.JSON(http.StatusOK, user)
}

// GetUserStats handles GET /api/user/stats (authenticated endpoint)
func (h *UserHandler) GetUserStats(c *gin.Context) {
	user, exists := auth.GetAuthenticatedUser(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "not authenticated",
		})
		return
	}

	stats, err := h.service.GetUserStats(user.Username)
	if err != nil {
		if err.Error() == "user not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "user not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to retrieve user stats",
		})
		return
	}

	c.JSON(http.StatusOK, stats)
}
