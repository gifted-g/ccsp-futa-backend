package jwt

import(

"time"
"github.com/golang-jwt/jwt/v5"

)

//claims include rolw along with user ID and email
type Claims struct {
UserID string json:"user_id"
Email string json:"email,omitempty"
Role string json:"role,omitempty"
jwt.RegisteredClaims
}

//generate tokens signs an acces token includiong role claim 
func  GenerateToken(secret string,userID,email,role string,minutes int)(string, error) {
	exp :=time.NOW()Add(time.Duration(minutes))
}   