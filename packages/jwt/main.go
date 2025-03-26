package jwt

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateAuthToken(userID, userName, userRole string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"name":    userName,
		"role":    userRole,
		"nbf":     time.Now().Add(time.Second * 3).Unix(),
		"exp":     time.Now().Add(time.Hour * 8).Unix(),
		"iat":     time.Now().Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))
	return tokenString, err
}
