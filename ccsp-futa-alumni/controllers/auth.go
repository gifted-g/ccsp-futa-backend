package controllers

import (
    "ccsp-futa-alumni/db"
    "ccsp-futa-alumni/models"
    "database/sql"
    "encoding/json"
   // "log"
    "net/http"
    "os"
    "time"

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

    _, err = db.DB.Exec(`INSERT INTO users (id, full_name, email, phone, password) VALUES ($1, $2, $3, $4, $5)`,
        id, input.FullName, input.Email, input.Phone, string(hashedPassword))

    if err != nil {
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


    err := db.DB.QueryRow(`SELECT id, full_name, email, phone, password, created_at FROM users WHERE email = $1`, input.Email).
        Scan(&user.ID, &user.FullName, &user.Email, &user.Phone, &user.Password, &user.CreatedAt)

    if err == sql.ErrNoRows {
        http.Error(w, "Email not found", http.StatusUnauthorized)
        return
    } else if err != nil {
        http.Error(w, "Server error", http.StatusInternalServerError)
        return
    }

    if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
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
