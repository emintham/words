package main

import (
	"log"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/words-api/words/internal/database"
	"github.com/words-api/words/internal/handlers"
)

func main() {
	// Initialize database
	db, err := database.InitDB("words.db")
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Create router
	router := gin.Default()

	// CORS middleware for web app
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Initialize handlers
	wordHandler := handlers.NewWordHandler(db)
	userHandler := handlers.NewUserHandler(db)
	vocabularyHandler := handlers.NewVocabularyHandler(db)
	reviewHandler := handlers.NewReviewHandler(db)

	// API routes
	api := router.Group("/api")
	{
		// Phase 1: Word lookup
		api.GET("/words/:word", wordHandler.GetWord)

		// Phase 2: User management
		api.POST("/users", userHandler.CreateUser)
		api.GET("/users/:username", userHandler.GetUser)
		api.GET("/users/:username/stats", userHandler.GetUserStats)

		// Phase 2: Vocabulary tracking
		api.POST("/users/:username/words/:word", vocabularyHandler.AddWord)
		api.GET("/users/:username/words", vocabularyHandler.GetUserWords)

		// Phase 2: Spaced repetition reviews
		api.GET("/users/:username/review", reviewHandler.GetDueWords)
		api.POST("/users/:username/review/:word", reviewHandler.SubmitReview)
		api.GET("/users/:username/review/:word/history", reviewHandler.GetReviewHistory)
	}

	// Get port from environment or default to 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting server on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
