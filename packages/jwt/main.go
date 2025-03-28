package jwt

import (
	"fmt"
	"os"
	"time"

	"github.com/application-ellas/ellas-backend/internal/domain/models"
	"github.com/golang-jwt/jwt/v5"
)

type CustomClaims struct {
	jwt.RegisteredClaims
	UserID string   `json:"user_id"`
	Name   string   `json:"name"`
	Roles  []string `json:"roles"`
}

func GenerateAuthToken(user models.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, CustomClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			NotBefore: jwt.NewNumericDate(time.Now().Add(time.Second * 3)),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 8)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		UserID: user.ID,
		Name:   user.Name,
		Roles:  user.GetUserRoles(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))
	return tokenString, err
}

func DecodeAuthToken(tokenString string) (*jwt.Token, *CustomClaims, error) {
	var claims CustomClaims
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("JWT_SECRET_KEY")), nil
	})

	return token, &claims, err
}
