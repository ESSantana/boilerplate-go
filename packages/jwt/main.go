package jwt

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/ESSantana/boilerplate-backend/packages/jwt/domain"
	"github.com/golang-jwt/jwt/v5"
)

func GenerateAuthToken(userID, userName, userRole string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, domain.CustomClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			NotBefore: jwt.NewNumericDate(time.Now().Add(time.Second * 3)),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 8)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		UserID: userID,
		Name:   userName,
		Role:   userRole,
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))
	return tokenString, err
}

func DecodeAuthToken(tokenString string) (domain.UserClaimData, error) {
	var claims domain.CustomClaims
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("JWT_SECRET_KEY")), nil
	})

	if err != nil {
		return domain.UserClaimData{}, fmt.Errorf("error parsing token: %v", err)
	}

	if !token.Valid {
		return domain.UserClaimData{}, fmt.Errorf("token is not valid")
	}

	return domain.UserClaimData{
		UserID: claims.UserID,
		Name:   claims.Name,
		Role:   claims.Role,
	}, nil
}

func ValidateUserRequestIssuer(request *http.Request, validateFunc func(...string) bool) error {
	tokenString := request.Header.Get("Authorization")
	if !strings.Contains(tokenString, "Bearer ") {
		return fmt.Errorf("authorization header missing Bearer token")
	}
	tokenString = tokenString[len("Bearer "):]

	claims, err := DecodeAuthToken(tokenString)
	if err != nil {
		return err
	}

	if !validateFunc(claims.Role, claims.UserID) {
		return fmt.Errorf("user validation failed")
	}

	return nil
}
