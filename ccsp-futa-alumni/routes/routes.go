package routes

import (
	"github.com/gin-gonic/gin"

	"ccsp-futa-alumni/handlers"
	"ccsp-futa-alumni/middleware"
)

// Public handler wrappers (so main.go can use them directly)
var RegisterHandler = handlers.RegisterHandler
var LoginHandler = handlers.LoginHandler
var LogoutHandler = handlers.LogoutHandler
var SendOTPHandler = handlers.SendOTPHandler
var VerifyOTPHandler = handlers.VerifyOTPHandler

func RegisterAuthRoutes(rg *gin.RouterGroup) {
	auth := rg.Group("/auth")
	{
		auth.GET("/me", middleware.RequireAuth(), handlers.MeHandler)
		// other auth-related
	}
}

func RegisterProfileRoutes(rg *gin.RouterGroup) {
	prof := rg.Group("/profile", middleware.RequireAuth())
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
	ev := rg.Group("/events", middleware.RequireAuth())
	{
		ev.GET("/", handlers.ListEventsHandler) // stub
	}
}

func RegisterMessageRoutes(rg *gin.RouterGroup) {
	msg := rg.Group("/messages", middleware.RequireAuth())
	{
		msg.POST("/channels", handlers.CreateChannelHandler)
		msg.GET("/channels/:channel_id", handlers.GetChannelHandler)
		msg.POST("/channels/:channel_id/members", handlers.AddChannelMemberHandler)
		msg.POST("/channels/:channel_id/messages", handlers.PostMessageHandler)
		msg.GET("/channels/:channel_id/messages", handlers.ListMessagesHandler)
	}
}

func RegisterContactRoutes(rg *gin.RouterGroup) {
	cont := rg.Group("/contact", middleware.RequireAuth())
	{
		cont.GET("/", handlers.ListContactsHandler) // stub
	}
}
