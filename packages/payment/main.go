package payment

import (
	"github.com/ESSantana/boilerplate-go/packages/log"
	"github.com/ESSantana/boilerplate-go/packages/payment/interfaces"
	"github.com/ESSantana/boilerplate-go/packages/payment/providers"
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
