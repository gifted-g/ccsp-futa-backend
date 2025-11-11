package middleware

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	// Changed UserID type to 'string' for consistency with jwt/jwt.go
	UserID string `json:"user_id"` 
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func AuthRequired(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		h := c.GetHeader("Authorization")
		if !strings.HasPrefix(h, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
			return
		}
		tokenStr := strings.TrimPrefix(h, "Bearer ")

		// Parse the token with the correct Claims struct
		t, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (interface{}, error) {
			// Security Best Practice: Check the signing method to prevent algorithm confusion attacks
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(secret), nil
		})

		if err != nil || !t.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}
		
		claims, ok := t.Claims.(*Claims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid claims format"})
			return
		}

		// The expiration check is usually done by jwt.ParseWithClaims, but we keep the logic clean:
		if claims.ExpiresAt != nil && time.Until(claims.ExpiresAt.Time) <= 0 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token expired"})
			return
		}
		
		// Set user ID and role in the context for downstream handlers
		c.Set("user_id", claims.UserID)
		c.Set("role", claims.Role)
		c.Next()
	}
}