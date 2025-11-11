package routes

import (
	"os" // <-- ADDED: Needed for os.Getenv
	
	"github.com/gin-gonic/gin"
	
	"ccsp-futa-alumni/handlers"
	"ccsp-futa-alumni/middleware"
)

// NOTE: Ensure JWT_SECRET is set in your environment variables.
// THIS IS THE SINGLE, AUTHORITATIVE DECLARATION OF jwtSecret for the entire 'routes' package.
var jwtSecret = os.Getenv("JWT_SECRET") // <-- UNCOMMENTED/FIXED

// Public handler wrappers (so main.go can use them directly)
var RegisterHandler = handlers.RegisterHandler
var LoginHandler = handlers.LoginHandler
var LogoutHandler = handlers.LogoutHandler
var SendOTPHandler = handlers.SendOTPHandler
var VerifyOTPHandler = handlers.VerifyOTPHandler

func RegisterAuthRoutes(rg *gin.RouterGroup) {
    auth := rg.Group("/auth")
    {
        auth.POST("/register", handlers.RegisterHandler)
        auth.POST("/login", handlers.LoginHandler)
        auth.POST("/logout", middleware.AuthRequired(jwtSecret), handlers.LogoutHandler)
        auth.POST("/send-otp", handlers.SendOTPHandler)
        auth.POST("/verify-otp", handlers.VerifyOTPHandler)
		auth.POST("/password/verify-otp", handlers.VerifyOTPHandler)

        auth.POST("/password/request-reset", handlers.RequestPasswordResetHandler)
        auth.POST("/password/reset", handlers.ResetPasswordHandler)
        auth.POST("/password/verify-otp", handlers.VerifyOTPHandler)

        auth.GET("/me", middleware.AuthRequired(jwtSecret), handlers.MeHandler)
    }
}


func RegisterProfileRoutes(rg *gin.RouterGroup) {
	// FIX: Using AuthRequired with the secret
	prof := rg.Group("/profile", middleware.AuthRequired(jwtSecret))
	{
		prof.GET("/:user_id", handlers.GetProfileHandler)
		prof.PUT("/:user_id", handlers.UpdateProfileHandler)
		prof.POST("/:user_id/avatar/request-upload", handlers.RequestAvatarUploadHandler)
		prof.POST("/:user_id/avatar/confirm", handlers.ConfirmAvatarHandler)
		// server-side upload endpoint (for dev convenience)
		prof.POST("/:user_id/avatar/upload", handlers.UploadAvatarHandler)
	}
}

func RegisterEventRoutes(rg *gin.RouterGroup) {
	// FIX: Using AuthRequired with the secret
	ev := rg.Group("/events", middleware.AuthRequired(jwtSecret))
	{
		// FIX: handlers.ListEventsHandler is now correctly referenced
		ev.GET("/", handlers.ListEventsHandler) // stub
	}
}

func RegisterMessageRoutes(rg *gin.RouterGroup) {
	// FIX: Using AuthRequired with the secret
	msg := rg.Group("/messages", middleware.AuthRequired(jwtSecret))
	{
		msg.POST("/channels", handlers.CreateChannelHandler)
		msg.GET("/channels/:channel_id", handlers.GetChannelHandler)
		msg.POST("/channels/:channel_id/members", handlers.AddChannelMemberHandler)
		msg.POST("/channels/:channel_id/messages", handlers.PostMessageHandler)
		msg.GET("/channels/:channel_id/messages", handlers.ListMessagesHandler)
	}
}

func RegisterContactRoutes(rg *gin.RouterGroup) {
	// FIX: Using AuthRequired with the secret
	cont := rg.Group("/contact", middleware.AuthRequired(jwtSecret))
	{
		// FIX: handlers.ListContactsHandler is now correctly referenced
		cont.GET("/", handlers.ListContactsHandler) // stub
	}
}

