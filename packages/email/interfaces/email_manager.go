package interfaces

import (
	"context"

	"github.com/ESSantana/boilerplate-backend/internal/domain/models"
)

type EmailManager interface {
	SendRecoverPasswordEmail(ctx context.Context, customer models.Customer) (err error)
}
