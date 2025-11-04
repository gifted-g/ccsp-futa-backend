package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"your_project/internal/auth"
)

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func LoginHandler(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// TODO: Check user in DB (hash check for password)
	// Assuming user found with ID=1
	userID := uint(1)

	// Generate JWT
	token, err := auth.GenerateJWT(userID, req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user":  gin.H{"id": userID, "email": req.Email},
	})
}
