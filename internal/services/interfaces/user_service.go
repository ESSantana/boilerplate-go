package interfaces

import (
	"context"

	"github.com/application-ellas/ellas-backend/internal/domain/models"
)

type UserService interface {
	GetUserByID(ctx context.Context, id string) (user models.User, err error)
	GetUserByExternalID(ctx context.Context, externalID string) (user models.User, err error)
	CreateUserIfNotExists(ctx context.Context, name, email, provider, externalID, profileImageURL string) (user models.User, err error)
}
