
package handlers

import (
    "context"
    "crypto/rand"
    "fmt"
    "github.com/dgrijalva/jwt-go"
    "github.com/gin-gonic/gin"
    "github.com/go-redis/redis/v8"
    "golang.org/x/crypto/bcrypt"
   /// "gorm.io/driver/sqlite"
    "gorm.io/gorm"
    "log"
    "net/http"
    "time"
) 

var (
    redisClient *redis.Client
    gormDB      *gorm.DB
    jwtKey      = []byte("secret_key") // Use env var in production
)


type Claims struct {
    Email string `json:"email"`
    jwt.StandardClaims
}

func generateOTP() string {
    b := make([]byte, 3)
    rand.Read(b)
    // Ensures a 6-digit number by masking and formatting
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
    
    // REMOVED DUPLICATE: otp := generateOTP()
    otp := generateOTP() 
    
    redisClient.Set(context.Background(), "reset_otp_"+req.Email, otp, 10*time.Minute)
    sendEmail(req.Email, "Password Reset OTP", "Your OTP: "+otp)
    c.JSON(http.StatusOK, gin.H{"message": "OTP sent"})
} // FIX: Missing closing brace added here

func VerifyResetOTP(c *gin.Context) {
    var req struct{ Email, OTP string }
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
        return
    }
    
    // FIX: Replaced 'ctx' with 'context.Background()' and removed duplicate line
    val, err := redisClient.Get(context.Background(), "reset_otp_"+req.Email).Result()
    
    if err != nil || val != req.OTP {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid/expired OTP"})
        return
    }

    // This logic was misplaced outside the function; it has been moved here.
    token := generateJWT(req.Email)
    redisClient.Set(context.Background(), "reset_token_"+req.Email, token, 10*time.Minute)
    c.JSON(http.StatusOK, gin.H{"reset_token": token})
} // FIX: Missing closing brace added here

func ResetPassword(c *gin.Context) {
    var req struct{ Email, ResetToken, NewPassword string }
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
        return
    }
    
    // FIX: Replaced 'ctx' with 'context.Background()' and removed duplicate line
    val, err := redisClient.Get(context.Background(), "reset_token_"+req.Email).Result()
    
    if err != nil || val != req.ResetToken {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid/expired token"})
        return
    }
    hashed, _ := hashPassword(req.NewPassword)
    gormDB.Model(&User{}).Where("email = ?", req.Email).Update("password_hash", hashed)
    redisClient.Del(context.Background(), "reset_token_"+req.Email)
    c.JSON(http.StatusOK, gin.H{"message": "Password reset successful"})
} // FIX: Missing closing brace added here

func SendEmailVerification(c *gin.Context) {
    var req struct{ Email string }
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
        return
    }
    
    // REMOVED DUPLICATE: otp := generateOTP()
    otp := generateOTP()
    
    redisClient.Set(context.Background(), "verify_email_"+req.Email, otp, 4*time.Minute)
    sendEmail(req.Email, "Email Verification OTP", "Your OTP: "+otp)
    c.JSON(http.StatusOK, gin.H{"message": "Verification OTP sent"})
} // FIX: Missing closing brace added here

func VerifyEmailOTP(c *gin.Context) {
    var req struct{ Email, OTP string }
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
        return
    }
    
    // FIX: Replaced 'ctx' with 'context.Background()' and removed duplicate line
    val, err := redisClient.Get(context.Background(), "verify_email_"+req.Email).Result()
    
    if err != nil || val != req.OTP {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid/expired OTP"})
        return
    }
    gormDB.Model(&User{}).Where("email = ?", req.Email).Update("email_verified", true)
    redisClient.Del(context.Background(), "verify_email_"+req.Email)
    c.JSON(http.StatusOK, gin.H{"message": "Email verified successfully"})
} // FIX: Missing closing brace added here

func generateJWT(email string) string {
    expirationTime := time.Now().Add(10 * time.Minute)
    claims := &Claims{Email: email, StandardClaims: jwt.StandardClaims{ExpiresAt: expirationTime.Unix()}}
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    tokenString, _ := token.SignedString(jwtKey)
    return tokenString
}

