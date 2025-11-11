package middleware

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

func RequireRole(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if v, ok := c.Get("role"); !ok || v.(string) != role {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error":"forbidden"})
			return
		}
		c.Next()
	}
}

func IsAdmin() gin.HandlerFunc {
    return func(c *gin.Context) {
        role, exists := c.Get("role")
        if !exists || role != "admin" {
            c.JSON(http.StatusForbidden, gin.H{"error": "admin privileges required"})
            c.Abort()
            return
        }
        c.Next()
    }
}
//(Assumes your JWT middleware in internal/auth/middleware.go sets c.Set("role", ...).)