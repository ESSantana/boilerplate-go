package middlewares

import (
	"net/http"
	"slices"
	"strings"

	"github.com/application-ellas/ella-backend/packages/jwt"
)

func AuthMiddleware(allowedRoles []string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
			authHeader := request.Header.Get("Authorization")
			if !strings.Contains(authHeader, "Bearer ") {
				response.WriteHeader(http.StatusUnauthorized)
				return
			}

			token := authHeader[len("Bearer "):]
			if token == "" {
				response.WriteHeader(http.StatusUnauthorized)
				return
			}

			jwt, customClaims, err := jwt.DecodeAuthToken(token)
			if err != nil {
				response.WriteHeader(http.StatusUnauthorized)
				return
			}

			if jwt.Valid {
				if slices.Contains(allowedRoles, customClaims.Role) {
					next.ServeHTTP(response, request)
				}
			}

			response.WriteHeader(http.StatusUnauthorized)
		})
	}

}
