package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"doodle-clone/internal/config"
	"doodle-clone/internal/database"
	"doodle-clone/internal/handlers"
	"doodle-clone/internal/middleware"
)

func main() {
	// Load configuration
	if err := config.Load(); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Set Gin mode
	if config.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}

	// Connect to database
	if err := database.Connect(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	// Run migrations
	if err := database.RunMigrations(); err != nil {
		log.Printf("Warning: migrations failed: %v", err)
		// Don't exit - maybe tables already exist
	}

	// Create router
	r := gin.Default()

	// Global middleware
	r.Use(middleware.CORS())
	r.Use(middleware.ErrorHandler())

	// Rate limiting for public endpoints
	r.Use(middleware.RateLimit(100, 60)) // 100 requests per minute

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"environment": config.AppConfig.Environment,
		})
	})

	// Create handlers
	authHandler := handlers.NewAuthHandler(database.Pool)
	pollHandler := handlers.NewPollHandler(database.Pool)
	voteHandler := handlers.NewVoteHandler(database.Pool)
	commentHandler := handlers.NewCommentHandler(database.Pool)
	exportHandler := handlers.NewExportHandler(database.Pool)

	// Google OAuth routes (without /api prefix for compatibility)
	google := r.Group("/auth")
	{
		google.GET("/google/login", authHandler.GoogleLogin)
		google.GET("/google/callback", authHandler.GoogleCallback)
	}

	// API routes
	api := r.Group("/api")
	{
		// Auth routes (public)
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.Refresh)
			auth.POST("/logout", authHandler.Logout)
			auth.GET("/me", middleware.Auth(), authHandler.GetMe)
			auth.PUT("/profile", middleware.Auth(), authHandler.UpdateProfile)
			auth.PUT("/password", middleware.Auth(), authHandler.ChangePassword)

			// Google OAuth (also under /api)
			auth.GET("/google/login", authHandler.GoogleLogin)
			auth.GET("/google/callback", authHandler.GoogleCallback)
		}

		// Public poll access
		api.GET("/polls", pollHandler.ListPolls)
		api.GET("/polls/:id", pollHandler.GetPoll)
		api.GET("/polls/:id/votes", voteHandler.GetVotes)
		api.GET("/polls/:id/comments", commentHandler.GetComments)

		// Exports (public)
		api.GET("/polls/:id/export/pdf", exportHandler.ExportPDF)
		api.GET("/polls/:id/export/ics", exportHandler.ExportICS)
		api.GET("/polls/:id/export/csv", exportHandler.ExportCSV)

		// Protected routes (require authentication)
		protected := api.Group("")
		protected.Use(middleware.Auth())
		{
			// Polls
			protected.POST("/polls", pollHandler.CreatePoll)
			protected.PUT("/polls/:id", pollHandler.UpdatePoll)
			protected.DELETE("/polls/:id", pollHandler.DeletePoll)
			protected.POST("/polls/:id/final", pollHandler.SetFinalDate)
			protected.POST("/polls/:id/dates", pollHandler.AddDateOption)

			// Votes
			protected.POST("/polls/:id/votes", voteHandler.CreateVote)
			protected.PUT("/polls/:id/votes/:voteId", voteHandler.UpdateVote)
			protected.DELETE("/polls/:id/votes/:voteId", voteHandler.DeleteVote)

			// Comments
			protected.POST("/polls/:id/comments", commentHandler.CreateComment)
			protected.PUT("/polls/:id/comments/:commentId", commentHandler.UpdateComment)
			protected.DELETE("/polls/:id/comments/:commentId", commentHandler.DeleteComment)

			// User dashboard
			protected.GET("/user/polls", pollHandler.GetUserPolls)
			protected.GET("/user/votes", voteHandler.GetUserVotes)
		}

		// Routes that support optional auth (can work with or without login)
		optionalAuth := api.Group("")
		optionalAuth.Use(middleware.OptionalAuth())
		{
			// Votes with optional auth (for anonymous voting)
			optionalAuth.POST("/polls/:id/vote", voteHandler.CreateVote)
		}
	}

	// Serve frontend in production
	if config.IsProduction() {
		// This would serve the built frontend files
		// For now, we'll just return a 404
		r.NoRoute(func(c *gin.Context) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
		})
	}

	// Start server
	port := config.AppConfig.Port
	if port == "" {
		port = "8080"
	}

	addr := ":" + port
	log.Printf("Starting server on %s", addr)
	log.Printf("Environment: %s", config.AppConfig.Environment)
	log.Printf("Frontend URL: %s", config.AppConfig.FrontendURL)

	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func init() {
	// Set up custom log format
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// Check for required env vars in production
	if os.Getenv("ENVIRONMENT") == "production" {
		required := []string{"DB_HOST", "DB_NAME", "DB_USER", "DB_PASSWORD", "JWT_SECRET"}
		for _, env := range required {
			if os.Getenv(env) == "" {
				log.Printf("Warning: Required environment variable %s is not set", env)
			}
		}
	}
}

// GetPort returns the port to listen on
func GetPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	return fmt.Sprintf(":%s", port)
}
