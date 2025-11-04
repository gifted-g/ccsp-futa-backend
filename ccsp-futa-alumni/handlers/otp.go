package handlers

import (
	"context"
	"crypto/rand"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"net/http"
//	"os"
	"time"
)  

var (
	ctx         = context.Background()
	redisClient *redis.Client
	db          *gorm.DB
	jwtKey      = []byte("secret_key") // Use env var in production
)

type User struct {
	ID            uint   `gorm:"primaryKey"`
	Email         string `gorm:"unique"`
	PasswordHash  string
	EmailVerified bool
}

type Claims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

func generateOTP() string {
	b := make([]byte, 3)
	rand.Read(b)
	return fmt.Sprintf("%06d", int(b[0])<<16|int(b[1])<<8|int(b[2])%1000000)
}

func sendEmail(to, subject, body string) {
	log.Printf("Sending email to %s: %s - %s\n", to, subject, body)
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

func RequestPasswordReset(c *gin.Context) {
	var req struct{ Email string }
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	otp := generateOTP()
	redisClient.Set(ctx, "reset_otp_"+req.Email, otp, 10*time.Minute)
	sendEmail(req.Email, "Password Reset OTP", "Your OTP: "+otp)
	c.JSON(http.StatusOK, gin.H{"message": "OTP sent"})
}

func VerifyResetOTP(c *gin.Context) {
	var req struct{ Email, OTP string }
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	val, err := redisClient.Get(ctx, "reset_otp_"+req.Email).Result()
	if err != nil || val != req.OTP {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid/expired OTP"})
		return
	}
	token := generateJWT(req.Email)
	redisClient.Set(ctx, "reset_token_"+req.Email, token, 10*time.Minute)
	c.JSON(http.StatusOK, gin.H{"reset_token": token})
}

func ResetPassword(c *gin.Context) {
	var req struct{ Email, ResetToken, NewPassword string }
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	val, err := redisClient.Get(ctx, "reset_token_"+req.Email).Result()
	if err != nil || val != req.ResetToken {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid/expired token"})
		return
	}
	hashed, _ := hashPassword(req.NewPassword)
	db.Model(&User{}).Where("email = ?", req.Email).Update("password_hash", hashed)
	redisClient.Del(ctx, "reset_token_"+req.Email)
	c.JSON(http.StatusOK, gin.H{"message": "Password reset successful"})
}

func SendEmailVerification(c *gin.Context) {
	var req struct{ Email string }
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	otp := generateOTP()
	redisClient.Set(ctx, "verify_email_"+req.Email, otp, 4*time.Minute)
	sendEmail(req.Email, "Email Verification OTP", "Your OTP: "+otp)
	c.JSON(http.StatusOK, gin.H{"message": "Verification OTP sent"})
}

func VerifyEmailOTP(c *gin.Context) {
	var req struct{ Email, OTP string }
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	val, err := redisClient.Get(ctx, "verify_email_"+req.Email).Result()
	if err != nil || val != req.OTP {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid/expired OTP"})
		return
	}
	db.Model(&User{}).Where("email = ?", req.Email).Update("email_verified", true)
	redisClient.Del(ctx, "verify_email_"+req.Email)
	c.JSON(http.StatusOK, gin.H{"message": "Email verified successfully"})
}

func generateJWT(email string) string {
	expirationTime := time.Now().Add(10 * time.Minute)
	claims := &Claims{Email: email, StandardClaims: jwt.StandardClaims{ExpiresAt: expirationTime.Unix()}}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString(jwtKey)
	return tokenString
}

func main() {
	dbTmp, _ := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	db = dbTmp
	db.AutoMigrate(&User{})
	redisClient = redis.NewClient(&redis.Options{Addr: "redis:6379"})
	r := gin.Default()

	r.POST("/auth/request-password-reset", RequestPasswordReset)
	r.POST("/auth/verify-otp", VerifyResetOTP)
	r.POST("/auth/reset-password", ResetPassword)
	r.POST("/auth/send-verification-otp", SendEmailVerification)
	r.POST("/auth/verify-email", VerifyEmailOTP)

	r.Run(":8080")
}
