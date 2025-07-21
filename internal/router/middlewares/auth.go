package middlewares

import (
	"net/http"
	"slices"
	"strings"

	"github.com/ESSantana/boilerplate-backend/internal/config"
	"github.com/ESSantana/boilerplate-backend/packages/jwt"
	"github.com/gofiber/fiber/v3"
)

func AuthMiddleware(cfg *config.Config, allowedRoles []string) fiber.Handler {
	return func(ctx fiber.Ctx) error {
		authHeader := ctx.Get("Authorization")
		if !strings.Contains(authHeader, "Bearer ") {
			ctx.Status(http.StatusUnauthorized)
			return nil
		}

		token := authHeader[len("Bearer "):]
		if token == "" {
			ctx.Status(http.StatusUnauthorized)
			return nil
		}

		userClaimData, err := jwt.DecodeAuthToken(cfg.JWT.SecretKey, token)
		if err != nil {
			ctx.Status(http.StatusUnauthorized)
			return nil
		}

		if slices.Contains(allowedRoles, userClaimData.Role) {
			return ctx.Next()
		}

		ctx.Status(http.StatusUnauthorized)
		return nil
	}

}
