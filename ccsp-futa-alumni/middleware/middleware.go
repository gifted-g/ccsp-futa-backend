package middleware

import (
    "time"

    "github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("your-secret-key") // Replace with secure secret from config

func GenerateToken(userID string) (string, error) {
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "sub": userID,
        "exp": time.Now().Add(time.Hour * 24).Unix(),
    })
    return token.SignedString(jwtSecret)
}