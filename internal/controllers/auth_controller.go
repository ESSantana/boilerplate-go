package controllers

import (
	"context"
	"errors"
	"os"
	"time"

	"net/http"

	"github.com/ESSantana/boilerplate-backend/internal/domain/constants"
	"github.com/ESSantana/boilerplate-backend/internal/domain/dto"
	"github.com/ESSantana/boilerplate-backend/internal/domain/models"
	svc_interfaces "github.com/ESSantana/boilerplate-backend/internal/services/interfaces"
	"github.com/ESSantana/boilerplate-backend/internal/utils"
	cache_interfaces "github.com/ESSantana/boilerplate-backend/packages/cache/interfaces"
	"github.com/ESSantana/boilerplate-backend/packages/email"
	email_interfaces "github.com/ESSantana/boilerplate-backend/packages/email/interfaces"
	"github.com/ESSantana/boilerplate-backend/packages/jwt"
	"github.com/ESSantana/boilerplate-backend/packages/log"
	"github.com/ESSantana/boilerplate-backend/packages/sso"
	sso_interfaces "github.com/ESSantana/boilerplate-backend/packages/sso/interfaces"
	"github.com/gofiber/fiber/v3"
)

type AuthController struct {
	logger         log.Logger
	serviceManager svc_interfaces.ServiceManager
	ssoManager     sso_interfaces.SSOManager
	cacheManager   cache_interfaces.CacheManager
	emailManager   email_interfaces.EmailManager
}

func NewAuthController(logger log.Logger, serviceManager svc_interfaces.ServiceManager, cacheManager cache_interfaces.CacheManager) AuthController {
	ssoManager := sso.NewSSOManager(
		cacheManager,
		sso.GoogleProvider{
			RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),
			ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
			ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		},
	)

	emailManager, err := email.NewEmailManager()
	if err != nil {
		logger.Errorf("failed to create email manager: %s", err.Error())
	}

	return AuthController{
		logger:         logger,
		serviceManager: serviceManager,
		ssoManager:     ssoManager,
		cacheManager:   cacheManager,
		emailManager:   emailManager,
	}
}

func (ctlr *AuthController) CustomerLogin(ctx fiber.Ctx) error {
	context, cancel := context.WithTimeout(ctx.Context(), constants.DefaultTimeout)
	defer cancel()

	loginRequest := utils.ReadBody[dto.LoginRequest](&ctx)
	if !loginRequest.IsValid() {
		utils.CreateResponse(&ctx, http.StatusBadRequest, errors.New("invalid login request"))
		return nil
	}
	ctlr.logger.Debugf("login request received: %v", loginRequest)

	customerService := ctlr.serviceManager.NewCustomerService()
	customer, err := customerService.GetCustomerLogin(context, loginRequest.Email, loginRequest.PasswordHash)
	if err != nil {
		utils.CreateResponse(&ctx, http.StatusInternalServerError, err)
		return nil
	}

	token, err := jwt.GenerateAuthToken(customer.ID, customer.Name, constants.RoleCustomer)
	if err != nil {
		ctlr.logger.Errorf("auth token error: %s", err.Error())
		utils.CreateResponse(&ctx, http.StatusInternalServerError, err)
		return nil
	}

	responseData := dto.LoginResponse{
		Token: token,
	}

	utils.CreateResponse(&ctx, http.StatusOK, nil, responseData)
	return nil
}

func (ctlr *AuthController) RecoverPassword(ctx fiber.Ctx) error {
	context, cancel := context.WithTimeout(ctx.Context(), constants.DefaultTimeout)
	defer cancel()

	recoverRequest := utils.ReadBody[dto.RecoverPasswordRequest](&ctx)
	if !recoverRequest.IsValid() {
		utils.CreateResponse(&ctx, http.StatusBadRequest, errors.New("invalid recover request"))
		return nil
	}
	ctlr.logger.Debugf("recover request received: %v", recoverRequest)

	customerService := ctlr.serviceManager.NewCustomerService()
	customer, err := customerService.GetCustomerByEmail(context, recoverRequest.Email)
	if err != nil {
		utils.CreateResponse(&ctx, http.StatusInternalServerError, err)
		return nil
	}

	responseData := dto.RecoverPasswordResponse{
		Message: "If you have an account with us, you will receive an email with instructions to reset your password.",
	}

	if customer.ID == "" {
		utils.CreateResponse(&ctx, http.StatusOK, nil, responseData)
		return nil
	}

	err = ctlr.emailManager.SendRecoverPasswordEmail(context, customer)
	if err != nil {
		utils.CreateResponse(&ctx, http.StatusInternalServerError, err)
		return nil
	}

	utils.CreateResponse(&ctx, http.StatusOK, nil, responseData)
	return nil
}

func (ctlr *AuthController) SSORequest(ctx fiber.Ctx) error {
	context, cancel := context.WithTimeout(ctx.Context(), constants.DefaultTimeout)
	defer cancel()

	providerName := ctx.Params("provider")
	provider, err := ctlr.ssoManager.GetProvider(providerName)
	if err != nil {
		utils.CreateResponse(&ctx, http.StatusInternalServerError, err)
		return nil
	}

	url, state := provider.GetSigninURL()
	err = ctlr.cacheManager.SetFlagWithExpiration(context, state, true, time.Minute*3)
	if err != nil {
		utils.CreateResponse(&ctx, http.StatusInternalServerError, err)
		return nil
	}

	responseBody := map[string]any{
		"redirect_url": url,
	}

	utils.CreateResponse(&ctx, http.StatusOK, nil, responseBody)
	return nil
}

func (ctlr *AuthController) SSOCallback(ctx fiber.Ctx) error {
	context, cancel := context.WithTimeout(ctx.Context(), constants.DefaultTimeout)
	defer cancel()

	providerName := ctx.Params("provider")
	provider, err := ctlr.ssoManager.GetProvider(providerName)
	if err != nil {
		utils.CreateResponse(&ctx, http.StatusInternalServerError, err)
		return nil
	}

	data, err := provider.GetUserData(ctx)
	if err != nil {
		utils.CreateResponse(&ctx, http.StatusInternalServerError, err)
		return nil
	}

	customer := models.Customer{
		Name:            data.Name,
		Email:           data.Email,
		ExternalID:      &data.ExternalID,
		ProfileImageURL: &data.ProfileImageURL,
	}

	userService := ctlr.serviceManager.NewCustomerService()
	customerCreated, err := userService.CreateCustomer(context, customer)
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

	url := os.Getenv("FRONTEND_URL") + "?token=" + token
	ctx.Redirect().To(url)
	ctx.Redirect().Status(http.StatusFound)

	return nil
}
