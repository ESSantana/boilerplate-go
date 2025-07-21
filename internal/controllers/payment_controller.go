package controllers

import (
	"context"
	"encoding/json"
	"errors"

	"net/http"

	"github.com/ESSantana/boilerplate-backend/internal/config"
	"github.com/ESSantana/boilerplate-backend/internal/domain/constants"
	"github.com/ESSantana/boilerplate-backend/internal/domain/dto"
	svc_interfaces "github.com/ESSantana/boilerplate-backend/internal/services/interfaces"
	"github.com/ESSantana/boilerplate-backend/internal/utils"
	"github.com/ESSantana/boilerplate-backend/packages/log"
	"github.com/ESSantana/boilerplate-backend/packages/payment"
	payment_interfaces "github.com/ESSantana/boilerplate-backend/packages/payment/interfaces"
	"github.com/gofiber/fiber/v3"
)

type PaymentController struct {
	cfg            *config.Config
	logger         log.Logger
	serviceManager svc_interfaces.ServiceManager
	paymentManager payment_interfaces.PaymentManager
}

func NewPaymentController(
	cfg *config.Config,
	logger log.Logger,
	serviceManager svc_interfaces.ServiceManager,
) PaymentController {
	paymentManager := payment.NewPaymentManager(cfg, logger)

	return PaymentController{
		cfg:            cfg,
		logger:         logger,
		serviceManager: serviceManager,
		paymentManager: paymentManager,
	}
}

func (ctlr *PaymentController) ExecutePayment(ctx fiber.Ctx) error {
	context, cancel := context.WithTimeout(ctx.Context(), constants.DefaultTimeout)
	defer cancel()

	requestBody := utils.ReadBody[dto.PaymentInfo](&ctx)
	if len(requestBody.Items) == 0 {
		utils.CreateResponse(&ctx, http.StatusBadRequest, errors.New("items cannot be empty"))
		return nil
	}
	ctlr.logger.Debugf("Payment request received: %v", requestBody)

	provider, err := ctlr.paymentManager.NewMercadoPagoProvider()
	if err != nil {
		utils.CreateResponse(&ctx, http.StatusInternalServerError, err)
		return nil
	}

	providerResponse, err := provider.ExecutePayment(context, requestBody)
	if err != nil {
		utils.CreateResponse(&ctx, http.StatusInternalServerError, &json.UnmarshalTypeError{})
		return nil
	}

	bytes, err := json.Marshal(*providerResponse)
	if err != nil {
		utils.CreateResponse(&ctx, http.StatusInternalServerError, err)
		return nil
	}

	var responseBody map[string]any
	err = json.Unmarshal(bytes, &responseBody)
	if err != nil {
		utils.CreateResponse(&ctx, http.StatusInternalServerError, err)
		return nil
	}

	utils.CreateResponse(&ctx, http.StatusOK, nil, responseBody)
	return nil
}

func (ctlr *PaymentController) PaymentWebhook(ctx fiber.Ctx) error {
	_, cancel := context.WithTimeout(ctx.Context(), constants.DefaultTimeout)
	defer cancel()

	requestBody := utils.ReadBody[map[string]any](&ctx)
	if requestBody == nil {
		return nil
	}

	ctlr.logger.Debugf("Payment webhook received: %v", requestBody)
	return nil
}
