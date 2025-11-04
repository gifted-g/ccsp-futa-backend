package routes

import (
	"github.com/gin-gonic/gin"
	"ccsp-futa-alumni/handlers"
)

func RegisterAuthRoutes(rg *gin.RouterGroup) {
	auth := rg.Group("/auth")
	{
		auth.POST("/register", handlers.RegisterHandler)
		auth.POST("/login", handlers.LoginHandler)
		auth.POST("/logout", handlers.LogoutHandler)
		auth.POST("/send-otp", handlers.SendOTPHandler)
		auth.POST("/verify-email", handlers.VerifyOTPHandler)
	}
}
