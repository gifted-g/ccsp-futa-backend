package handlers

import (
    "context"
    "crypto/rand"
    "database/sql"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "os"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/go-redis/redis/v8"
    "golang.org/x/crypto/bcrypt"
    _ "github.com/lib/pq"
)

var (
    db  *sql.DB
    rdb *redis.Client
    ctx = context.Background()
)

func init() {
    var err error
    db, err = sql.Open("postgres", os.Getenv("DATABASE_URL"))
    if err != nil {
        log.Fatal(err)
    }

    rdb = redis.NewClient(&redis.Options{
        Addr:     os.Getenv("REDIS_ADDR"),
        Password: os.Getenv("REDIS_PASSWORD"),
        DB:       0,
    })
}

type User struct {
    Email    string `json:"email"`
    Password string `json:"password"`
}

// Generate a 6-digit OTP
func generateOTPReset() string {
    b := make([]byte, 3)
    rand.Read(b)
    return fmt.Sprintf("%06d", int(b[0])<<16|int(b[1])<<8|int(b[2]))[:6]
}

// Rate limiting key
func rateLimitKey(email string) string {
    return "rate_limit:" + email
}

// Check if user is rate-limited
func isRateLimited(email string) bool {
    key := rateLimitKey(email)
    count, err := rdb.Get(ctx, key).Int()
    if err != nil && err != redis.Nil {
        return true
    }
    if count >= 5 {
        return true
    }
    rdb.Incr(ctx, key)
    rdb.Expire(ctx, key, 1*time.Hour)
    return false
}

// Internal logic: Request password reset
func requestPasswordReset(w http.ResponseWriter, r *http.Request) {
    var user User
    json.NewDecoder(r.Body).Decode(&user)

    if isRateLimited(user.Email) {
        http.Error(w, "Too many requests. Try again later.", http.StatusTooManyRequests)
        return
    }

    otp := generateOTPReset()
    rdb.Set(ctx, "otp:"+user.Email, otp, 10*time.Minute)

    fmt.Printf("Sending OTP %s to email %s\n", otp, user.Email)
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("OTP sent to email"))
}

// Internal logic: Verify OTP
func verifyOTP(w http.ResponseWriter, r *http.Request) {
    var user User
    json.NewDecoder(r.Body).Decode(&user)

    storedOtp, err := rdb.Get(ctx, "otp:"+user.Email).Result()
    if err == redis.Nil {
        http.Error(w, "OTP expired or not found", http.StatusBadRequest)
        return
    }
    if user.Password != storedOtp {
        http.Error(w, "Invalid OTP", http.StatusUnauthorized)
        return
    }

    resetToken := generateOTPReset()
    rdb.Set(ctx, "reset:"+user.Email, resetToken, 15*time.Minute)

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{"reset_token": resetToken})
}

// Internal logic: Reset password
func resetPassword(w http.ResponseWriter, r *http.Request) {
    var req struct {
        Email       string `json:"email"`
        NewPassword string `json:"new_password"`
        ResetToken  string `json:"reset_token"`
    }
    json.NewDecoder(r.Body).Decode(&req)

    storedToken, err := rdb.Get(ctx, "reset:"+req.Email).Result()
    if err == redis.Nil {
        http.Error(w, "Reset token expired or invalid", http.StatusBadRequest)
        return
    }
    if req.ResetToken != storedToken {
        http.Error(w, "Invalid reset token", http.StatusUnauthorized)
        return
    }

    hash, _ := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
    _, err = db.Exec("UPDATE users SET password=$1 WHERE email=$2", string(hash), req.Email)
    if err != nil {
        http.Error(w, "Error updating password", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
    w.Write([]byte("Password reset successful"))
}

//
// Gin-compatible handlers
//

func RequestPasswordResetHandler(c *gin.Context) {
    requestPasswordReset(c.Writer, c.Request)
}


func ResetPasswordHandler(c *gin.Context) {
    resetPassword(c.Writer, c.Request)
}
