package main

import (
	"log"
	"os"

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

	// Initialize handlers
	wordHandler := handlers.NewWordHandler(db)

	// API routes
	api := router.Group("/api")
	{
		api.GET("/words/:word", wordHandler.GetWord)
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
