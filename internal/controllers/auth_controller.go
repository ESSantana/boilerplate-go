package controllers

import (
	"context"
	"os"
	"time"

	"net/http"

	"github.com/application-ellas/ellas-backend/internal/domain/constants"
	"github.com/application-ellas/ellas-backend/internal/domain/models"
	svc_interfaces "github.com/application-ellas/ellas-backend/internal/services/interfaces"
	"github.com/application-ellas/ellas-backend/internal/utils"
	cache_interfaces "github.com/application-ellas/ellas-backend/packages/cache/interfaces"
	"github.com/application-ellas/ellas-backend/packages/jwt"
	"github.com/application-ellas/ellas-backend/packages/log"
	"github.com/application-ellas/ellas-backend/packages/sso"
	sso_interfaces "github.com/application-ellas/ellas-backend/packages/sso/interfaces"
	"github.com/go-chi/chi/v5"
)

type AuthController struct {
	logger         log.Logger
	serviceManager svc_interfaces.ServiceManager
	ssoManager     sso_interfaces.SSOManager
	cacheManager   cache_interfaces.CacheManager
}

func NewAuthController(logger log.Logger, serviceManager svc_interfaces.ServiceManager, cacheManager cache_interfaces.CacheManager) AuthController {
	ssoManager := sso.NewSSOManager(
		cacheManager,
		sso.GoogleProvider{
			RedirectURL:  os.Getenv("REDIRECT_URL"),
			ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
			ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		},
	)

	return AuthController{
		logger:         logger,
		serviceManager: serviceManager,
		ssoManager:     ssoManager,
		cacheManager:   cacheManager,
	}
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
