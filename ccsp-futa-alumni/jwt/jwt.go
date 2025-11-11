package jwt

import (
	"fmt" // FIX: Added the fmt package for fmt.Errorf
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// claims include role along with user ID and email
type Claims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email,omitempty"`
	Role   string `json:"role,omitempty"`
	jwt.RegisteredClaims
}

// GenerateToken signs an access token including role claim
func GenerateToken(secret string, userID, email, role string, minutes int) (string, error) {
	expTime := time.Now().Add(time.Duration(minutes) * time.Minute)

	claims := &Claims{
		UserID: userID,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			// A common practice is to use "sub" (subject) for the user ID
			Subject: userID,
			// Expiration Time
			ExpiresAt: jwt.NewNumericDate(expTime),
			// Issued At
			IssuedAt: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	
	// Sign the token using the provided secret
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		// Uses the imported fmt package
		return "", fmt.Errorf("failed to sign token: %w", err) 
	}

	return tokenString, nil
}

// NOTE: You will also likely need a 'VerifyToken' function to complete the package.