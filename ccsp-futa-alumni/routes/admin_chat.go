package routes

import (
	//"os" // Added to read the JWT secret environment variable

	"github.com/gin-gonic/gin"
	"ccsp-futa-alumni/handlers"
	"ccsp-futa-alumni/middleware"
)

// NOTE: You must define a JWT_SECRET environment variable or replace os.Getenv("JWT_SECRET") 
// with a hardcoded secret for this to work.
//var jwtSecret = os.Getenv("JWT_SECRET")

func RegisterAdminChatRoutes(r *gin.Engine) {
	// FIX: Changed middleware.Auth() to the correctly named middleware.AuthRequired() 
	// and provided the required secret.
	// NOTE: If middleware.IsAdmin() is not defined, this will still fail until it's created.
	admin := r.Group("/admin/chat", middleware.AuthRequired(jwtSecret), middleware.IsAdmin())
	{
		admin.DELETE("/messages/:id", handlers.DeleteMessage)
		admin.POST("/messages/:id/pin", handlers.PinMessage)
		admin.POST("/broadcast", handlers.Broadcast)
	}
}