package providers

import (
	"context"
	"encoding/json"
	"os"

	"github.com/ESSantana/boilerplate-backend/internal/domain/dto"
	"github.com/ESSantana/boilerplate-backend/packages/log"
	"github.com/ESSantana/boilerplate-backend/packages/payment/interfaces"
	"github.com/mercadopago/sdk-go/pkg/config"
	"github.com/mercadopago/sdk-go/pkg/preference"
)

const (
	DefaultCurrency = "BRL" //real brasileiro
)

type mercadoPagoPaymentProvider struct {
	logger         log.Logger
	mercadoPagoCfg *config.Config
}

func NewMercadoPagoProvider(logger log.Logger) (interfaces.PaymentProvider, error) {
	cfg, err := config.New(os.Getenv("MERCADO_PAGO_TOKEN"))
	if err != nil {
		return nil, err
	}

	return &mercadoPagoPaymentProvider{
		mercadoPagoCfg: cfg,
		logger:         logger,
	}, nil
}

func (provider *mercadoPagoPaymentProvider) ExecutePayment(ctx context.Context, paymentInfo dto.PaymentInfo) (*preference.Response, error) {
	items := make([]preference.ItemRequest, len(paymentInfo.Items))
	for i, item := range paymentInfo.Items {
		items[i] = preference.ItemRequest{
			ID:         item.ID,
			Title:      item.Title,
			UnitPrice:  item.UnitPrice,
			Quantity:   item.Quantity,
			CurrencyID: DefaultCurrency,
		}
	}
	preferenceRequest := preference.Request{
		BackURLs: &preference.BackURLsRequest{
			Success: "https://localhost:8080/success",
			Failure: "https://localhost:8080/failure",
			Pending: "https://localhost:8080/pending",
		},
		Payer: &preference.PayerRequest{
			Phone: &preference.PhoneRequest{
				AreaCode: paymentInfo.GetPhoneAreaCode(),
				Number:   paymentInfo.GetPhoneNumber(),
			},
			Name:    paymentInfo.GetFirstName(),
			Surname: paymentInfo.GetLastName(),
			Email:   paymentInfo.CustomerEmail,
		},
		Items:      items,
		AutoReturn: "all",
		PaymentMethods: &preference.PaymentMethodsRequest{
			ExcludedPaymentTypes: []preference.ExcludedPaymentTypeRequest{
				{ID: "ticket"},
				{ID: "boleto"},
			},
		},
	}

	data, _ := json.Marshal(preferenceRequest)

	provider.logger.Debugf("Creating preference with request: %s", string(data))

	preferenceClient := preference.NewClient(provider.mercadoPagoCfg)
	response, err := preferenceClient.Create(ctx, preferenceRequest)
	return response, err
}
