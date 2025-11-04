package routes

import (
	"net/http"

	"ccsp-futa-alumni/middleware"
	"ccsp-futa-alumni/ws"

	"github.com/gin-gonic/gin"
)

func WsHandler(c *gin.Context) {
	// Accept token either in header or query param
	auth := c.GetHeader("Authorization")
	if auth == "" {
		// try query
		auth = c.Query("token")
		if auth == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing token"}); return
		}
	} else {
		// remove Bearer prefix
		if len(auth) > 7 && auth[:7] == "Bearer " {
			auth = auth[7:]
		}
	}
	tok, err := middleware.ParseToken(auth)
	if err != nil || !tok.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"}); return
	}
	claims := tok.Claims.(map[string]interface{})
	sub, _ := claims["sub"].(string)

	// Hijack connection
	ws.ServeWS(c.Writer, c.Request, sub)
}
