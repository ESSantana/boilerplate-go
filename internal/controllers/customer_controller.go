package controllers

import (
	"context"
	"encoding/json"
	"errors"

	"net/http"

	"github.com/ESSantana/boilerplate-backend/internal/domain/constants"
	"github.com/ESSantana/boilerplate-backend/internal/domain/dto"
	"github.com/ESSantana/boilerplate-backend/internal/domain/models"
	svc_interfaces "github.com/ESSantana/boilerplate-backend/internal/services/interfaces"
	"github.com/ESSantana/boilerplate-backend/internal/utils"
	cache_interfaces "github.com/ESSantana/boilerplate-backend/packages/cache/interfaces"
	"github.com/ESSantana/boilerplate-backend/packages/jwt"
	"github.com/ESSantana/boilerplate-backend/packages/log"
	"github.com/go-chi/chi/v5"
)

type CustomerController struct {
	logger         log.Logger
	serviceManager svc_interfaces.ServiceManager
	cacheManager   cache_interfaces.CacheManager
}

func NewCustomerController(logger log.Logger, serviceManager svc_interfaces.ServiceManager, cacheManager cache_interfaces.CacheManager) CustomerController {
	return CustomerController{
		logger:         logger,
		serviceManager: serviceManager,
		cacheManager:   cacheManager,
	}
}

func (ctlr *CustomerController) GetCustomerById(response http.ResponseWriter, request *http.Request) {
	context, cancel := context.WithTimeout(request.Context(), constants.DefaultTimeout)
	defer cancel()

	customerId := chi.URLParam(request, "id")
	ctlr.logger.Debugf("customer_id: %v", customerId)

	err := jwt.ValidateUserRequestIssuer(request, utils.CreateUserValidation(customerId))
	if err != nil {
		utils.CreateResponse(&response, http.StatusForbidden, err)
		return
	}

	customerService := ctlr.serviceManager.NewCustomerService()
	customer, err := customerService.GetCustomerById(context, customerId)
	if err != nil {
		utils.CreateResponse(&response, http.StatusInternalServerError, err)
		return
	}
	if customer.ID == "" {
		utils.CreateResponse(&response, http.StatusNoContent, nil, "customer not found")
		return
	}

	utils.CreateResponse(&response, http.StatusOK, nil, customer)
}

func (ctlr *CustomerController) GetAllCustomers(response http.ResponseWriter, request *http.Request) {
	context, cancel := context.WithTimeout(request.Context(), constants.DefaultTimeout)
	defer cancel()

	customerService := ctlr.serviceManager.NewCustomerService()
	customers, err := customerService.GetAllCustomers(context)
	if err != nil {
		utils.CreateResponse(&response, http.StatusInternalServerError, err)
		return
	}
	if len(customers) == 0 {
		utils.CreateResponse(&response, http.StatusNoContent, nil, "any customer found")
		return
	}

	utils.CreateResponse(&response, http.StatusOK, nil, customers)
}

func (ctlr *CustomerController) Create(response http.ResponseWriter, request *http.Request) {
	context, cancel := context.WithTimeout(request.Context(), constants.DefaultTimeout)
	defer cancel()

	var customer models.Customer
	err := json.NewDecoder(request.Body).Decode(&customer)
	if err != nil {
		utils.CreateResponse(&response, http.StatusBadRequest, errors.New("payload format is invalid"))
		return
	}
	ctlr.logger.Debugf("customer received: %v", customer)

	customerService := ctlr.serviceManager.NewCustomerService()
	customerCreated, err := customerService.CreateCustomer(context, customer)
	if err != nil {
		utils.CreateResponse(&response, http.StatusInternalServerError, err)
		return
	}

	token, err := jwt.GenerateAuthToken(customerCreated.ID, customerCreated.Name, constants.RoleCustomer)
	if err != nil {
		ctlr.logger.Errorf("auth token error: %s", err.Error())
		utils.CreateResponse(&response, http.StatusInternalServerError, err)
		return
	}

	responseData := dto.LoginResponse{
		Token: token,
	}

	utils.CreateResponse(&response, http.StatusCreated, nil, responseData)
}

func (ctlr *CustomerController) Update(response http.ResponseWriter, request *http.Request) {
	context, cancel := context.WithTimeout(request.Context(), constants.DefaultTimeout)
	defer cancel()

	var customer models.Customer
	err := json.NewDecoder(request.Body).Decode(&customer)
	if err != nil {
		utils.CreateResponse(&response, http.StatusBadRequest, errors.New("payload format is invalid"))
		return
	}
	ctlr.logger.Debugf("customer received: %v", customer)

	err = jwt.ValidateUserRequestIssuer(request, utils.CreateUserValidation(customer.ID))
	if err != nil {
		utils.CreateResponse(&response, http.StatusForbidden, err)
		return
	}

	customerService := ctlr.serviceManager.NewCustomerService()
	err = customerService.UpdateCustomer(context, customer)
	if err != nil {
		utils.CreateResponse(&response, http.StatusInternalServerError, err)
		return
	}
	utils.CreateResponse(&response, http.StatusOK, nil, "customer updated successfully")
}

func (ctlr *CustomerController) SoftDelete(response http.ResponseWriter, request *http.Request) {
	context, cancel := context.WithTimeout(request.Context(), constants.DefaultTimeout)
	defer cancel()

	var customer models.Customer
	err := json.NewDecoder(request.Body).Decode(&customer)
	if err != nil {
		utils.CreateResponse(&response, http.StatusBadRequest, errors.New("payload format is invalid"))
		return
	}
	ctlr.logger.Debugf("customer received: %v", customer)

	err = jwt.ValidateUserRequestIssuer(request, utils.CreateUserValidation(customer.ID))
	if err != nil {
		utils.CreateResponse(&response, http.StatusForbidden, err)
		return
	}

	customerService := ctlr.serviceManager.NewCustomerService()
	err = customerService.SoftDeleteCustomer(context, customer.ID)
	if err != nil {
		utils.CreateResponse(&response, http.StatusInternalServerError, err)
		return
	}
	utils.CreateResponse(&response, http.StatusOK, nil, "customer deleted successfully")
}
