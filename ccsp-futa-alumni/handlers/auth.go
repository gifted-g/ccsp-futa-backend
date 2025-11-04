package handlers

import (
	"net/http"
	"time"

	"ccsp-futa-alumni/db"
	"ccsp-futa-alumni/models"
	"ccsp-futa-alumni/middleware"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type registerReq struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

func RegisterHandler(c *gin.Context) {
	var req registerReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// hash password
	h, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	u := models.User{
		ID:           uuid.New(),
		Email:        req.Email,
		PasswordHash: string(h),
	}
	if err := db.DB.Create(&u).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email taken or invalid"})
		return
	}
	// create empty profile
	p := models.Profile{
		UserID:      u.ID,
		DisplayName: "",
	}
	db.DB.Create(&p)

	token, _ := middleware.GenerateToken(u.ID.String())
	c.JSON(http.StatusCreated, gin.H{
		"user":  gin.H{"id": u.ID, "email": u.Email},
		"token": token,
	})
}

type loginReq struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func LoginHandler(c *gin.Context) {
	var req loginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var u models.User
	if err := db.DB.Where("email = ?", req.Email).First(&u).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	token, _ := middleware.GenerateToken(u.ID.String())
	c.JSON(http.StatusOK, gin.H{
		"user":  gin.H{"id": u.ID, "email": u.Email},
		"token": token,
	})
}

func LogoutHandler(c *gin.Context) {
	// For stateless JWT, logout can be client-side (delete token). Optionally you can implement a token blacklist.
	c.JSON(http.StatusOK, gin.H{"msg": "logged out"})
}

// OTP stubs (you can wire email provider)
func SendOTPHandler(c *gin.Context) {
	// stub: accept email, "send" otp
	type R struct{ Email string `json:"email" binding:"required,email"` }
	var r R
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}); return
	}
	// TODO: persist OTP & send via email
	c.JSON(http.StatusOK, gin.H{"msg": "otp sent (stub)", "email": r.Email, "expires_at": time.Now().Add(5 * time.Minute)})
}

func VerifyOTPHandler(c *gin.Context) {
	// stub: verify OTP from body
	c.JSON(http.StatusOK, gin.H{"msg": "otp verified (stub)"})
}

func MeHandler(c *gin.Context) {
	sub := c.GetString("sub")
	var u models.User
	if err := db.DB.Preload("Profile").Where("id = ?", sub).First(&u).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"}); return
	}
	c.JSON(http.StatusOK, gin.H{"user": u})
}
