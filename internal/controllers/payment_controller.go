package controllers

import (
	"context"
	"encoding/json"

	"net/http"

	"github.com/ESSantana/boilerplate-backend/internal/domain/constants"
	"github.com/ESSantana/boilerplate-backend/internal/domain/dto"
	svc_interfaces "github.com/ESSantana/boilerplate-backend/internal/services/interfaces"
	"github.com/ESSantana/boilerplate-backend/internal/utils"
	"github.com/ESSantana/boilerplate-backend/packages/log"
	"github.com/ESSantana/boilerplate-backend/packages/payment"
	payment_interfaces "github.com/ESSantana/boilerplate-backend/packages/payment/interfaces"
)

type PaymentController struct {
	logger         log.Logger
	serviceManager svc_interfaces.ServiceManager
	paymentManager payment_interfaces.PaymentManager
}

func NewPaymentController(logger log.Logger, serviceManager svc_interfaces.ServiceManager) PaymentController {
	paymentManager := payment.NewPaymentManager(logger)

	return PaymentController{
		logger:         logger,
		serviceManager: serviceManager,
		paymentManager: paymentManager,
	}
}

func (ctlr *PaymentController) ExecutePayment(response http.ResponseWriter, request *http.Request) {
	ctxWithTimeout, cancel := context.WithTimeout(request.Context(), constants.DefaultTimeout)
	defer cancel()

	requestBody := utils.ReadBody[dto.PaymentInfo](request, response)
	if len(requestBody.Items) == 0 {
		return
	}

	ctlr.logger.Debugf("Payment request received: %v", requestBody)

	provider, err := ctlr.paymentManager.NewMercadoPagoProvider()
	if err != nil {
		utils.CreateResponse(&response, http.StatusInternalServerError, err)
		return
	}

	providerResponse, err := provider.ExecutePayment(ctxWithTimeout, requestBody)
	if err != nil {
		utils.CreateResponse(&response, http.StatusInternalServerError, &json.UnmarshalTypeError{})
		return
	}

	bytes, err := json.Marshal(*providerResponse)
	if err != nil {
		utils.CreateResponse(&response, http.StatusInternalServerError, err)
		return
	}

	var responseBody map[string]interface{}
	err = json.Unmarshal(bytes, &responseBody)
	if err != nil {
		utils.CreateResponse(&response, http.StatusInternalServerError, err)
		return
	}

	utils.CreateResponse(&response, http.StatusOK, nil, responseBody)
}

func (ctlr *PaymentController) PaymentWebhook(response http.ResponseWriter, request *http.Request) {
	_, cancel := context.WithTimeout(request.Context(), constants.DefaultTimeout)
	defer cancel()

	requestBody := utils.ReadBody[map[string]interface{}](request, response)
	if requestBody == nil {
		return
	}

	ctlr.logger.Debugf("Payment webhook received: %v", requestBody)
}
