package payment

import (
	"github.com/ESSantana/boilerplate-backend/internal/config"
	"github.com/ESSantana/boilerplate-backend/packages/payment/interfaces"
	"github.com/ESSantana/boilerplate-backend/packages/payment/providers"
)

type PaymentManager struct {
  cfg   *config.Config
}

func NewPaymentManager(cfg *config.Config) interfaces.PaymentManager {
	return &PaymentManager{
    cfg:   cfg,
	}
}

func (pm *PaymentManager) NewMercadoPagoProvider() (interfaces.PaymentProvider, error) {
	return providers.NewMercadoPagoProvider( pm.cfg.MercadoPago.Token)
}
