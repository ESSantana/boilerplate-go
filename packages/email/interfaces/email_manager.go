package interfaces

import (
	"context"

	"github.com/application-ellas/ella-backend/internal/domain/models"
)

type EmailManager interface {
	SendRecoverPasswordEmail(ctx context.Context, customer models.Customer) (err error)
}
