package interfaces

import (
	"context"
)

type ServiceProviderService interface {
	PromoteUserToServiceProvider(ctx context.Context, userID string) (err error)
}
