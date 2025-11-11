package routes

import (
	"net/http"
	//"os" // Added back to use os.Getenv
	"strings" 

	"ccsp-futa-alumni/middleware"
	"ccsp-futa-alumni/ws"

	"github.com/gin-gonic/gin"
)

// NOTE: Variable is reinstated to fix the 'undefined' error for jwtSecret. 
// If this declaration causes a 'redeclared' error later, it means this
// variable needs to be moved to a single file (or a config package).
//var jwtSecret = os.Getenv("JWT_SECRET")

func WsHandler(c *gin.Context) {
	// 1. Extract Token from Header or Query
	tokenStr := ""
	auth := c.GetHeader("Authorization")
	if auth == "" {
		tokenStr = c.Query("token")
	} else if strings.HasPrefix(auth, "Bearer ") {
		tokenStr = strings.TrimPrefix(auth, "Bearer ")
	} else {
		tokenStr = auth
	}

	if tokenStr == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
		return
	}

	// 2. Parse and Validate Token using the reusable middleware function
	// This line assumes middleware.ParseToken is exported (starts with a capital P)
	// and defined in the ccsp-futa-alumni/middleware package.
	tok, err := middleware.ParseToken(tokenStr, jwtSecret) 
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token: " + err.Error()})
		return
	}

	// 3. Extract Claims
	// This line assumes the Claims struct is also exported from the middleware package.
	claims, ok := tok.Claims.(*middleware.Claims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid claims format"})
		return
	}
	sub := claims.UserID // Use UserID from the custom Claims struct

	// 4. Hijack connection
	ws.ServeWS(c.Writer, c.Request, sub)
}