package interfaces

import (
	"context"

	"github.com/application-ellas/ellas-backend/internal/domain/models"
)

type UserRepository interface {
	GetUserById(ctx context.Context, id string) (user models.User, err error)
	GetUserByExternalId(ctx context.Context, externalID string) (user models.User, err error)
	GetUserRoleByUserIdAndRole(ctx context.Context, userID, role string) (userRole models.UserRole, err error)
	AppendRoleToUser(ctx context.Context, userID, role string) (err error)
	CreateUser(ctx context.Context, user models.User) (err error)
}
