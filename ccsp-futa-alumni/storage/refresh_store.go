package storage

import (
    "time"
    "math/rand"
)

var refreshTokens = make(map[string]struct {
    UserID string
    Expiry time.Time
})

func SaveRefreshToken(userID string) string {
    token := RandString(32)
    refreshTokens[token] = struct {
        UserID string
        Expiry time.Time
    }{UserID: userID, Expiry: time.Now().Add(24 * time.Hour)}
    return token
}

func ValidateRefreshToken(token string) (string, bool) {
    data, exists := refreshTokens[token]
    if !exists || time.Now().After(data.Expiry) {
        return "", false
    }
    return data.UserID, true
}

func RandString(n int) string {
    letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
    b := make([]rune, n)
    for i := range b {
        b[i] = letters[rand.Intn(len(letters))]
    }
    return string(b)
}
