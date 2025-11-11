package controllers

import (
    "net/http"
    "os"
    "time"

    "ccsp-futa-alumni/config"
    "ccsp-futa-alumni/models"
    
    // The imported package is now used to fix the logic
    "golang.org/x/crypto/bcrypt" 
    "github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte(os.Getenv("JWT_SECRET"))

func authControllerLogin(c *gin.Context) {
    var input struct {
        Email    string `json:"email"`
        Password string `json:"password"`
    }
    c.BindJSON(&input)

    var user models.User
    db := config.DB
    if err := db.Where("email = ?", input.Email).First(&user).Error; err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
        return
    }

    // --- FIX APPLIED HERE ---
    // 1. Replaced models.CheckPasswordHash with bcrypt.CompareHashAndPassword.
    // 2. Changed user.Password to the correct struct field: user.PasswordHash.
    if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
        return
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "sub": user.ID, // Added sub as a standard JWT claim
        "email": user.Email,
        "role":  user.Role,
        "exp":   time.Now().Add(time.Hour * 72).Unix(),
    })

    tokenString, _ := token.SignedString(jwtKey)
    c.JSON(http.StatusOK, gin.H{"token": tokenString})
}