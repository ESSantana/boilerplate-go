package controllers

import (
	"context"

	"net/http"

	"github.com/ESSantana/boilerplate-backend/internal/domain/constants"
	"github.com/ESSantana/boilerplate-backend/internal/domain/dto"
	"github.com/ESSantana/boilerplate-backend/internal/domain/models"
	svc_interfaces "github.com/ESSantana/boilerplate-backend/internal/services/interfaces"
	"github.com/ESSantana/boilerplate-backend/internal/utils"
	cache_interfaces "github.com/ESSantana/boilerplate-backend/packages/cache/interfaces"
	"github.com/ESSantana/boilerplate-backend/packages/jwt"
	"github.com/ESSantana/boilerplate-backend/packages/log"
	"github.com/gofiber/fiber/v3"
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

func (ctlr *CustomerController) GetCustomerById(ctx fiber.Ctx) error {
	context, cancel := context.WithTimeout(ctx.Context(), constants.DefaultTimeout)
	defer cancel()

	customerId := ctx.Params("id")
	ctlr.logger.Debugf("customer_id: %v", customerId)

	err := jwt.ValidateUserRequestIssuer(ctx, utils.CreateUserValidation(customerId))
	if err != nil {
		utils.CreateResponse(&ctx, http.StatusForbidden, err)
		return nil
	}

	customerService := ctlr.serviceManager.NewCustomerService()
	customer, err := customerService.GetCustomerById(context, customerId)
	if err != nil {
		utils.CreateResponse(&ctx, http.StatusInternalServerError, err)
		return nil
	}

	if customer.ID == "" {
		utils.CreateResponse(&ctx, http.StatusNoContent, nil, "customer not found")
		return nil
	}

	utils.CreateResponse(&ctx, http.StatusOK, nil, customer)
	return nil
}

func (ctlr *CustomerController) GetAllCustomers(ctx fiber.Ctx) error {
	context, cancel := context.WithTimeout(ctx.Context(), constants.DefaultTimeout)
	defer cancel()

	customerService := ctlr.serviceManager.NewCustomerService()
	customers, err := customerService.GetAllCustomers(context)
	if err != nil {
		utils.CreateResponse(&ctx, http.StatusInternalServerError, err)
		return nil
	}
  
	if len(customers) == 0 {
		utils.CreateResponse(&ctx, http.StatusNoContent, nil, "any customer found")
		return nil
	}

	utils.CreateResponse(&ctx, http.StatusOK, nil, customers)
	return nil
}

func (ctlr *CustomerController) Create(ctx fiber.Ctx) error {
	context, cancel := context.WithTimeout(ctx.Context(), constants.DefaultTimeout)
	defer cancel()

	customer := utils.ReadBody[models.Customer](&ctx)
	ctlr.logger.Debugf("customer received: %v", customer)

	customerService := ctlr.serviceManager.NewCustomerService()
	customerCreated, err := customerService.CreateCustomer(context, customer)
	if err != nil {
		utils.CreateResponse(&ctx, http.StatusInternalServerError, err)
		return nil
	}

	token, err := jwt.GenerateAuthToken(customerCreated.ID, customerCreated.Name, constants.RoleCustomer)
	if err != nil {
		ctlr.logger.Errorf("auth token error: %s", err.Error())
		utils.CreateResponse(&ctx, http.StatusInternalServerError, err)
		return nil
	}

	responseData := dto.LoginResponse{
		Token: token,
	}

	utils.CreateResponse(&ctx, http.StatusCreated, nil, responseData)
	return nil
}

func (ctlr *CustomerController) Update(ctx fiber.Ctx) error {
	context, cancel := context.WithTimeout(ctx.Context(), constants.DefaultTimeout)
	defer cancel()

	customer := utils.ReadBody[models.Customer](&ctx)
	ctlr.logger.Debugf("customer received: %v", customer)

	err := jwt.ValidateUserRequestIssuer(ctx, utils.CreateUserValidation(customer.ID))
	if err != nil {
		utils.CreateResponse(&ctx, http.StatusForbidden, err)
		return nil
	}

	customerService := ctlr.serviceManager.NewCustomerService()
	err = customerService.UpdateCustomer(context, customer)
	if err != nil {
		utils.CreateResponse(&ctx, http.StatusInternalServerError, err)
		return nil
	}

	utils.CreateResponse(&ctx, http.StatusOK, nil, "customer updated successfully")
	return nil
}

func (ctlr *CustomerController) SoftDelete(ctx fiber.Ctx) error {
	context, cancel := context.WithTimeout(ctx.Context(), constants.DefaultTimeout)
	defer cancel()

	customer := utils.ReadBody[models.Customer](&ctx)
	ctlr.logger.Debugf("customer received: %v", customer)

	err := jwt.ValidateUserRequestIssuer(ctx, utils.CreateUserValidation(customer.ID))
	if err != nil {
		utils.CreateResponse(&ctx, http.StatusForbidden, err)
		return nil
	}

	customerService := ctlr.serviceManager.NewCustomerService()
	err = customerService.SoftDeleteCustomer(context, customer.ID)
	if err != nil {
		utils.CreateResponse(&ctx, http.StatusInternalServerError, err)
		return nil
	}

	utils.CreateResponse(&ctx, http.StatusOK, nil, "customer deleted successfully")
	return nil
}
