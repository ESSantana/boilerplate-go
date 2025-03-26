package controllers

import (
	"context"
	"os"
	"time"

	"net/http"

	"github.com/ESSantana/boilerplate-go/internal/domain/constants"
	svc_interfaces "github.com/ESSantana/boilerplate-go/internal/services/interfaces"
	"github.com/ESSantana/boilerplate-go/internal/utils"
	cache_interfaces "github.com/ESSantana/boilerplate-go/packages/cache/interfaces"
	"github.com/ESSantana/boilerplate-go/packages/jwt"
	"github.com/ESSantana/boilerplate-go/packages/log"
	"github.com/ESSantana/boilerplate-go/packages/sso"
	sso_interfaces "github.com/ESSantana/boilerplate-go/packages/sso/interfaces"
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
		body := map[string]interface{}{
			"error":   true,
			"message": err.Error(),
		}
		utils.CreateResponse(&response, http.StatusInternalServerError, body)
		return
	}

	url, state := provider.GetSigninURL()
	err = ctlr.cacheManager.SetFlagWithExpiration(request.Context(), state, true, time.Minute*3)
	if err != nil {
		body := map[string]interface{}{
			"error":   true,
			"message": err.Error(),
		}
		utils.CreateResponse(&response, http.StatusInternalServerError, body)
		return
	}

	responseBody := map[string]interface{}{
		"redirect_url": url,
	}

	utils.CreateResponse(&response, http.StatusOK, responseBody)
}

func (ctlr *AuthController) SSOCallback(response http.ResponseWriter, request *http.Request) {
	ctx, cancel := context.WithTimeout(request.Context(), constants.DefaultTimeout)
	defer cancel()

	providerName := chi.URLParam(request, "provider")
	provider, err := ctlr.ssoManager.GetProvider(providerName)
	if err != nil {
		body := map[string]interface{}{
			"error":   true,
			"message": err.Error(),
		}
		utils.CreateResponse(&response, http.StatusInternalServerError, body)
		return
	}

	data, err := provider.GetUserData(request)
	if err != nil {
		body := map[string]interface{}{
			"error":   true,
			"message": err.Error(),
		}
		utils.CreateResponse(&response, http.StatusInternalServerError, body)
		return
	}

	userService := ctlr.serviceManager.NewUserService()
	user, err := userService.CreateUserIfNotExists(ctx, data.Name, data.Email, providerName, data.ExternalID, data.ProfileImageURL)
	if err != nil {
		body := map[string]interface{}{
			"error":   true,
			"message": err.Error(),
		}
		utils.CreateResponse(&response, http.StatusInternalServerError, body)
		return
	}

	token, err := jwt.GenerateAuthToken(user.ID, user.Name, "user")
	if err != nil {
		ctlr.logger.Errorf("auth token error: %s", err.Error())
		body := map[string]interface{}{
			"error":   true,
			"message": "error generating auth token",
		}
		utils.CreateResponse(&response, http.StatusInternalServerError, body)
		return
	}

	url := os.Getenv("FRONTEND_URL") + "?token=" + token

	http.Redirect(response, request, url, http.StatusSeeOther)
}
