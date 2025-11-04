package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID int    `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func AuthRequired(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		h := c.GetHeader("Authorization")
		if !strings.HasPrefix(h, "Bearer ") { c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error":"missing token"}); return }
		tokenStr := strings.TrimPrefix(h, "Bearer ")
		t, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token)(interface{}, error){ return []byte(secret), nil })
		if err != nil || !t.Valid { c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error":"invalid token"}); return }
		claims := t.Claims.(*Claims)
		if claims.ExpiresAt != nil && time.Until(claims.ExpiresAt.Time) <= 0 { c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error":"token expired"}); return }
		c.Set("user_id", claims.UserID)
		c.Set("role", claims.Role)
		c.Next()

		

}
	}

}

