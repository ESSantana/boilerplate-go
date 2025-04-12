package controllers

import (
	"context"
	"encoding/json"
	"os"
	"time"

	"net/http"

	"github.com/application-ellas/ella-backend/internal/domain/constants"
	"github.com/application-ellas/ella-backend/internal/domain/dto"
	"github.com/application-ellas/ella-backend/internal/domain/errors"
	"github.com/application-ellas/ella-backend/internal/domain/models"
	svc_interfaces "github.com/application-ellas/ella-backend/internal/services/interfaces"
	"github.com/application-ellas/ella-backend/internal/utils"
	cache_interfaces "github.com/application-ellas/ella-backend/packages/cache/interfaces"
	"github.com/application-ellas/ella-backend/packages/email"
	email_interfaces "github.com/application-ellas/ella-backend/packages/email/interfaces"
	"github.com/application-ellas/ella-backend/packages/jwt"
	"github.com/application-ellas/ella-backend/packages/log"
	"github.com/application-ellas/ella-backend/packages/sso"
	sso_interfaces "github.com/application-ellas/ella-backend/packages/sso/interfaces"
	"github.com/go-chi/chi/v5"
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

func (ctlr *AuthController) CustomerLogin(response http.ResponseWriter, request *http.Request) {
	context, cancel := context.WithTimeout(request.Context(), constants.DefaultTimeout)
	defer cancel()

	var loginRequest dto.LoginRequest
	err := json.NewDecoder(request.Body).Decode(&loginRequest)
	if err != nil {
		utils.CreateResponse(&response, http.StatusBadRequest, errors.New("payload format is invalid"))
		return
	}
	ctlr.logger.Debugf("login request received: %v", loginRequest)

	customerService := ctlr.serviceManager.NewCustomerService()
	customer, err := customerService.GetCustomerLogin(context, loginRequest.Email, loginRequest.PasswordHash)
	if err != nil {
		utils.CreateResponse(&response, http.StatusInternalServerError, err)
		return
	}

	token, err := jwt.GenerateAuthToken(customer.ID, customer.Name, constants.RoleCustomer)
	if err != nil {
		ctlr.logger.Errorf("auth token error: %s", err.Error())
		utils.CreateResponse(&response, http.StatusInternalServerError, err)
		return
	}

	responseData := dto.LoginResponse{
		Token: token,
	}

	utils.CreateResponse(&response, http.StatusOK, nil, responseData)
}

func (ctlr *AuthController) RecoverPassword(response http.ResponseWriter, request *http.Request) {
	context, cancel := context.WithTimeout(request.Context(), constants.DefaultTimeout)
	defer cancel()

	var recoverRequest dto.RecoverPasswordRequest
	err := json.NewDecoder(request.Body).Decode(&recoverRequest)
	if err != nil {
		utils.CreateResponse(&response, http.StatusBadRequest, errors.New("payload format is invalid"))
		return
	}
	ctlr.logger.Debugf("recover request received: %v", recoverRequest)

	customerService := ctlr.serviceManager.NewCustomerService()
	customer, err := customerService.GetCustomerByEmail(context, recoverRequest.Email)
	if err != nil {
		utils.CreateResponse(&response, http.StatusInternalServerError, err)
		return
	}

	responseData := dto.RecoverPasswordResponse{
		Message: "If you have an account with us, you will receive an email with instructions to reset your password.",
	}

	if customer.ID == "" {
		utils.CreateResponse(&response, http.StatusOK, nil, responseData)
		return
	}

	err = ctlr.emailManager.SendRecoverPasswordEmail(context, customer)
	if err != nil {
		utils.CreateResponse(&response, http.StatusInternalServerError, err)
		return
	}

	utils.CreateResponse(&response, http.StatusOK, nil, responseData)
}

func (ctlr *AuthController) SSORequest(response http.ResponseWriter, request *http.Request) {
	_, cancel := context.WithTimeout(request.Context(), constants.DefaultTimeout)
	defer cancel()

	providerName := chi.URLParam(request, "provider")
	provider, err := ctlr.ssoManager.GetProvider(providerName)
	if err != nil {
		utils.CreateResponse(&response, http.StatusInternalServerError, err)
		return
	}

	url, state := provider.GetSigninURL()
	err = ctlr.cacheManager.SetFlagWithExpiration(request.Context(), state, true, time.Minute*3)
	if err != nil {
		utils.CreateResponse(&response, http.StatusInternalServerError, err)
		return
	}

	responseBody := map[string]interface{}{
		"redirect_url": url,
	}

	utils.CreateResponse(&response, http.StatusOK, nil, responseBody)
}

func (ctlr *AuthController) SSOCallback(response http.ResponseWriter, request *http.Request) {
	ctx, cancel := context.WithTimeout(request.Context(), constants.DefaultTimeout)
	defer cancel()

	providerName := chi.URLParam(request, "provider")
	provider, err := ctlr.ssoManager.GetProvider(providerName)
	if err != nil {
		utils.CreateResponse(&response, http.StatusInternalServerError, err)
		return
	}

	data, err := provider.GetUserData(request)
	if err != nil {
		utils.CreateResponse(&response, http.StatusInternalServerError, err)
		return
	}

	customer := models.Customer{
		Name:            data.Name,
		Email:           data.Email,
		ExternalID:      &data.ExternalID,
		ProfileImageURL: &data.ProfileImageURL,
	}

	userService := ctlr.serviceManager.NewCustomerService()
	customerCreated, err := userService.CreateCustomer(ctx, customer)
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

	url := os.Getenv("FRONTEND_URL") + "?token=" + token
	http.Redirect(response, request, url, http.StatusSeeOther)
}
