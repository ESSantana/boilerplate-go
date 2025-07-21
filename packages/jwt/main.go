package jwt

import (
	"fmt"
	"strings"
	"time"

	"github.com/ESSantana/boilerplate-backend/packages/jwt/domain"
	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
)

func GenerateAuthToken(jwtSecret, userID, userName, userRole string) (string, error) {
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

	tokenString, err := token.SignedString([]byte(jwtSecret))
	return tokenString, err
}

func DecodeAuthToken(jwtSecret, tokenString string) (domain.UserClaimData, error) {
	var claims domain.CustomClaims
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtSecret), nil
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

func ValidateUserRequestIssuer(ctx fiber.Ctx, jwtSecret string, validateFunc func(...string) bool) error {
	tokenString := ctx.Get("Authorization")
	if !strings.Contains(tokenString, "Bearer ") {
		return fmt.Errorf("authorization header missing Bearer token")
	}
	tokenString = tokenString[len("Bearer "):]

	claims, err := DecodeAuthToken(jwtSecret, tokenString)
	if err != nil {
		return err
	}

	if !validateFunc(claims.Role, claims.UserID) {
		return fmt.Errorf("user validation failed")
	}

	return nil
}
