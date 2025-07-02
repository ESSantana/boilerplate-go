package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/ESSantana/boilerplate-backend/internal/domain/constants"
	"github.com/ESSantana/boilerplate-backend/internal/domain/dto"
	"github.com/ESSantana/boilerplate-backend/internal/domain/errors"
	"github.com/gofiber/fiber/v3"
)

// CreateResponse creates a standardized HTTP response with the given status code, error, and data.
func CreateResponse(ctx *fiber.Ctx, statusCode int, responseErr error, data ...any) {
	var body = dto.HttpResponse{}
	if responseErr != nil {
		body.Error = true
		body.Message = responseErr.Error()
	}

	if len(data) > 0 {
		body.Data = data[0]
	}

	if _, ok := responseErr.(*errors.ValidationError); ok {
		statusCode = http.StatusUnprocessableEntity
	}

	if _, ok := responseErr.(*errors.NotFoundError); ok {
		statusCode = http.StatusNotFound
	}

	if _, ok := responseErr.(*errors.ForbiddenError); ok {
		statusCode = http.StatusForbidden
	}

	c := *ctx

	c.Status(statusCode)
	c.Set("Content-Type", "application/json")
	err := c.JSON(body)
	if err != nil {
		c.Status(http.StatusInternalServerError)
	}

	ctx = &c
}

// ReadBody reads the request body and unmarshals it into the provided type T.
func ReadBody[T any](ctx *fiber.Ctx) (output T) {
	var bodyRequest T
	json.Unmarshal((*ctx).Body(), &bodyRequest)
	return bodyRequest
}

// CreateUserValidation returns a function that validates user creation based on the role and user ID.
func CreateUserValidation(expectedID string) func(params ...string) bool {
	return func(params ...string) bool {
		if len(params) < 2 {
			fmt.Println("params length is less than 2")
			return false
		}
		role := strings.TrimSpace(params[0])
		userID := strings.TrimSpace(params[1])

		if role == constants.RoleAdmin || role == constants.RoleManager {
			return true
		}

		return userID == expectedID
	}
}
