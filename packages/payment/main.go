package payment

import (
	"github.com/application-ellas/ella-backend/packages/log"
	"github.com/application-ellas/ella-backend/packages/payment/interfaces"
	"github.com/application-ellas/ella-backend/packages/payment/providers"
)

type PaymentManager struct {
	logger log.Logger
}

func NewPaymentManager(logger log.Logger) interfaces.PaymentManager {
	return &PaymentManager{
		logger: logger,
	}
}

func (pm *PaymentManager) NewMercadoPagoProvider() (interfaces.PaymentProvider, error) {
	return providers.NewMercadoPagoProvider(pm.logger)
}
