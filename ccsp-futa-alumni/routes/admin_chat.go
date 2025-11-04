package routes

import (
    "github.com/gin-gonic/gin"
    "yourapp/handlers"
    "yourapp/middleware"
)

func RegisterAdminChatRoutes(r *gin.Engine) {
    admin := r.Group("/admin/chat", middleware.Auth(), middleware.IsAdmin())
    {
        admin.DELETE("/messages/:id", handlers.DeleteMessage)
        admin.POST("/messages/:id/pin", handlers.PinMessage)
        admin.POST("/broadcast", handlers.Broadcast)
    }
}
