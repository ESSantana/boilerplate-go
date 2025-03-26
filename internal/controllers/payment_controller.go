package controllers

import (
	"context"
	"encoding/json"

	"net/http"

	"github.com/ESSantana/boilerplate-go/internal/domain/constants"
	"github.com/ESSantana/boilerplate-go/internal/domain/dto"
	svc_interfaces "github.com/ESSantana/boilerplate-go/internal/services/interfaces"
	"github.com/ESSantana/boilerplate-go/internal/utils"
	"github.com/ESSantana/boilerplate-go/packages/log"
	"github.com/ESSantana/boilerplate-go/packages/payment"
	payment_interfaces "github.com/ESSantana/boilerplate-go/packages/payment/interfaces"
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
		body := map[string]interface{}{
			"error":   true,
			"message": err.Error(),
		}
		utils.CreateResponse(&response, http.StatusInternalServerError, body)
		return
	}

	providerResponse, err := provider.ExecutePayment(ctxWithTimeout, requestBody)
	if err != nil {
		body := map[string]interface{}{
			"error":   true,
			"message": err.Error(),
		}
		utils.CreateResponse(&response, http.StatusInternalServerError, body)
		return
	}

	bytes, err := json.Marshal(*providerResponse)
	if err != nil {
		body := map[string]interface{}{
			"error":   true,
			"message": err.Error(),
		}
		utils.CreateResponse(&response, http.StatusInternalServerError, body)
		return
	}

	var responseBody map[string]interface{}
	err = json.Unmarshal(bytes, &responseBody)
	if err != nil {
		body := map[string]interface{}{
			"error":   true,
			"message": err.Error(),
		}
		utils.CreateResponse(&response, http.StatusInternalServerError, body)
		return
	}

	utils.CreateResponse(&response, http.StatusOK, responseBody)
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
