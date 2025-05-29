package interfaces

import (
	"context"

	"github.com/ESSantana/boilerplate-backend/internal/domain/dto"
	"github.com/mercadopago/sdk-go/pkg/preference"
)

type PaymentProvider interface {
	ExecutePayment(ctx context.Context, paymentInfo dto.PaymentInfo) (*preference.Response, error)
}
