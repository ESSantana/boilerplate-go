package payment

import (
	"github.com/ESSantana/boilerplate-backend/internal/config"
	"github.com/ESSantana/boilerplate-backend/packages/log"
	"github.com/ESSantana/boilerplate-backend/packages/payment/interfaces"
	"github.com/ESSantana/boilerplate-backend/packages/payment/providers"
)

type PaymentManager struct {
  cfg   *config.Config
	logger log.Logger
}

func NewPaymentManager(cfg *config.Config, logger log.Logger) interfaces.PaymentManager {
	return &PaymentManager{
    cfg:   cfg,
		logger: logger,
	}
}

func (pm *PaymentManager) NewMercadoPagoProvider() (interfaces.PaymentProvider, error) {
	return providers.NewMercadoPagoProvider(pm.logger, pm.cfg.MercadoPago.Token)
}
