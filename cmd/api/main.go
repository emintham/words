package main

import (
	"flag"
	"log"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/words-api/words/internal/auth"
	"github.com/words-api/words/internal/database"
	"github.com/words-api/words/internal/handlers"
)

func main() {
	// Parse command-line flags
	portFlag := flag.String("port", "", "Port to run the server on")
	flag.Parse()

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

	// Initialize session store
	sessionStore := auth.NewSessionStore()

	// Initialize handlers
	wordHandler := handlers.NewWordHandler(db)
	userHandler := handlers.NewUserHandler(db)
	vocabularyHandler := handlers.NewVocabularyHandler(db)
	reviewHandler := handlers.NewReviewHandler(db)
	authHandler := handlers.NewAuthHandler(db, sessionStore)

	// API routes
	api := router.Group("/api")
	{
		// Public routes
		// Phase 1: Word lookup (public)
		api.GET("/words/:word", wordHandler.GetWord)

		// Authentication routes (public)
		api.POST("/auth/login", authHandler.Login)
		api.POST("/auth/logout", authHandler.Logout)

		// User registration (public)
		api.POST("/users", userHandler.CreateUser)

		// Protected routes (require authentication)
		protected := api.Group("")
		protected.Use(auth.AuthMiddleware(sessionStore))
		{
			// User management
			protected.GET("/user", userHandler.GetUser)
			protected.GET("/user/stats", userHandler.GetUserStats)
			protected.GET("/auth/me", authHandler.GetCurrentUser)

			// Vocabulary tracking
			protected.POST("/words/:word", vocabularyHandler.AddWord)
			protected.GET("/words", vocabularyHandler.GetUserWords)

			// Spaced repetition reviews
			protected.GET("/review", reviewHandler.GetDueWords)
			protected.POST("/review/:word", reviewHandler.SubmitReview)
			protected.GET("/review/:word/history", reviewHandler.GetReviewHistory)
		}
	}

	// Determine port: command-line flag > environment variable > default
	port := *portFlag
	if port == "" {
		port = os.Getenv("PORT")
	}
	if port == "" {
		port = "9090"
	}

	log.Printf("Starting server on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
