package interfaces

import (
	"github.com/ESSantana/boilerplate-backend/internal/domain/dto"
	"github.com/gofiber/fiber/v3"
)

type SSOProvider interface {
	GetSigninURL() (signinURL, userState string)
	GetUserData(ctx fiber.Ctx) (data dto.UserDataSSO, err error)
}
