package interfaces

import (
	"context"

	"github.com/ESSantana/boilerplate-go/internal/domain/models"
)

type UserRepository interface {
	GetUserById(ctx context.Context, id string) (user models.User, err error)
	GetUserByExternalId(ctx context.Context, externalID string) (user models.User, err error)
	CreateUser(ctx context.Context, user models.User) (err error)
}
