package routes

import (
	"github.com/gin-gonic/gin"
	"ccsp-futa-alumni/handlers"
)

// Renamed function to RegisterAuthSubRoutes to avoid conflict with a similarly named 
// function in routes/routes.go (as indicated by the error details).
func RegisterAuthSubRoutes(rg *gin.RouterGroup) {
	auth := rg.Group("/auth")
	{
		auth.POST("/register", handlers.RegisterHandler)
		auth.POST("/login", handlers.LoginHandler)
		auth.POST("/logout", handlers.LogoutHandler)
		auth.POST("/send-otp", handlers.SendOTPHandler)
		auth.POST("/verify-email", handlers.VerifyOTPHandler)
	}
}