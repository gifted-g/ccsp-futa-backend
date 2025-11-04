package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetUserProfile(c *gin.Context) {
	userID, _ := c.Get("userID")
	email, _ := c.Get("email")

	c.JSON(http.StatusOK, gin.H{
		"message": "Profile fetched successfully",
		"user": gin.H{
			"id":    userID,
			"email": email,
		},
	})
}
