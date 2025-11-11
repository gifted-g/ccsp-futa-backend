package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"

	"ccsp-futa-alumni/db"
)

func main() {
	// Initialize DB connection (db.Init does not return an error)
	db.Init()

	// Create Gin router
	router := gin.Default()

	// Health check / root route
	router.GET("/", func(c *gin.Context) {
		c.String(200, "Welcome to CCSP Alumni Backend! 🚀")
	})

	// Public routes (placeholder implementations to avoid undefined references)
	router.POST("/register", func(c *gin.Context) {
		c.JSON(501, gin.H{"error": "not implemented"})
	})
	router.POST("/login", func(c *gin.Context) {
		c.JSON(501, gin.H{"error": "not implemented"})
	})
	router.POST("/logout", func(c *gin.Context) {
		c.JSON(501, gin.H{"error": "not implemented"})
	})
	router.POST("/send-otp", func(c *gin.Context) {
		c.JSON(501, gin.H{"error": "not implemented"})
	})
	router.POST("/verify-email", func(c *gin.Context) {
		c.JSON(501, gin.H{"error": "not implemented"})
	})
	router.Static("/uploads", "./uploads")

	// Protected API routes (grouped under /api/v1)
	//api := router.Group("/api/v1")
	//{
		// Intentionally left empty: register routes from the routes package
		// when those functions are implemented to avoid build errors.
	//}
	api := router.Group("/api/v1")
	api.GET("/", func(c *gin.Context) {
	c.JSON(200, gin.H{"message": "API v1 Root"})
	})




	
	// Get port from env or fallback to 3000
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Printf("🚀 Server starting on port %s...\n", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("❌ Could not start server: %v", err)
	}
	if err := db.Init(); err != nil {
		log.Fatalf("❌ Failed to initialize database: %v", err)
	}
}
	

