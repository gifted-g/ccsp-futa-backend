package controllers

import (
    "ccsp-futa-alumni/db" // Assumed GORM connection is here
    "ccsp-futa-alumni/models"
   // "database/sql"
    "encoding/json"
    "net/http"
    "os"
    "time"
    "errors" 
    "gorm.io/gorm"
    "github.com/golang-jwt/jwt/v5"
    "golang.org/x/crypto/bcrypt"
    "github.com/google/uuid"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

func Signup(w http.ResponseWriter, r *http.Request) {
    var input struct {
        FullName string `json:"full_name"`
        Email    string `json:"email"`
        Phone    string `json:"phone"`
        Password string `json:"password"`
    }

    if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
        http.Error(w, "Invalid input", http.StatusBadRequest)
        return
    }

    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
    if err != nil {
        http.Error(w, "Error processing password", http.StatusInternalServerError)
        return
    }

    id := uuid.New().String()

    // FIX 1: Use GORM's Exec which returns a *gorm.DB object. Check the .Error property.
    result := db.DB.Exec(`INSERT INTO users (id, full_name, email, phone, password_hash) VALUES ($1, $2, $3, $4, $5)`,
        id, input.FullName, input.Email, input.Phone, string(hashedPassword))

    // Check the error property on the GORM result object
    if result.Error != nil {
        http.Error(w, "Email or phone already exists", http.StatusConflict)
        return
    }

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(map[string]string{"message": "Signup successful"})
}

func Login(w http.ResponseWriter, r *http.Request) {
    var input struct {
        Email    string `json:"email"`
        Password string `json:"password"`
    }

    if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
        http.Error(w, "Invalid input", http.StatusBadRequest)
        return
    }
 
    var user models.User

    // FIX 1: Use GORM's Raw().Scan() to execute the query and map result to the 'user' struct.
    // FIX 2: Assuming 'user' struct now contains 'FullName' and 'PasswordHash'
    err := db.DB.Raw(`SELECT id, full_name, email, phone, password_hash, created_at FROM users WHERE email = ?`, input.Email).
        Scan(&user).Error 

    if errors.Is(err, gorm.ErrRecordNotFound) { // Check for GORM's 'no rows' error
        http.Error(w, "Email not found", http.StatusUnauthorized)
        return
    } else if err != nil {
        http.Error(w, "Server error", http.StatusInternalServerError)
        return
    }

    if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
        http.Error(w, "Incorrect password", http.StatusUnauthorized)
        return
    }

    // Generate JWT
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "user_id": user.ID,
        "exp":     time.Now().Add(72 * time.Hour).Unix(),
    })

    tokenString, err := token.SignedString(jwtSecret)
    if err != nil {
        http.Error(w, "Token generation failed", http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
}