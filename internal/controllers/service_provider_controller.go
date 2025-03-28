package controllers

import (
	"context"
	"encoding/json"
	"io"

	"net/http"

	"github.com/application-ellas/ellas-backend/internal/domain/constants"
	"github.com/application-ellas/ellas-backend/internal/domain/dto"
	svc_interfaces "github.com/application-ellas/ellas-backend/internal/services/interfaces"
	"github.com/application-ellas/ellas-backend/internal/utils"
	cache_interfaces "github.com/application-ellas/ellas-backend/packages/cache/interfaces"
	"github.com/application-ellas/ellas-backend/packages/log"
)

type ServiceProviderController struct {
	logger         log.Logger
	serviceManager svc_interfaces.ServiceManager
	cacheManager   cache_interfaces.CacheManager
}

func NewServiceProviderController(logger log.Logger, serviceManager svc_interfaces.ServiceManager, cacheManager cache_interfaces.CacheManager) ServiceProviderController {
	return ServiceProviderController{
		logger:         logger,
		serviceManager: serviceManager,
		cacheManager:   cacheManager,
	}
}

func (ctlr *ServiceProviderController) PromoteUserToServiceProvider(response http.ResponseWriter, request *http.Request) {
	context, cancel := context.WithTimeout(request.Context(), constants.DefaultTimeout)
	defer cancel()

	body, err := io.ReadAll(request.Body)
	if err != nil {
		utils.CreateResponse(&response, http.StatusBadRequest, map[string]interface{}{
			"error":   true,
			"message": "Invalid request body",
		})
		return
	}

	var creationRequest dto.ServiceProviderCreationRequest
	err = json.Unmarshal(body, &creationRequest)
	if err != nil {
		utils.CreateResponse(&response, http.StatusBadRequest, map[string]interface{}{
			"error":   true,
			"message": "Invalid request body",
		})
		return
	}

	ctlr.logger.Debugf("Promote user to service provider request received: %v", creationRequest)

	serviceProviderService := ctlr.serviceManager.NewServiceProviderService()
	err = serviceProviderService.PromoteUserToServiceProvider(context, creationRequest.UserID)
	if err != nil {
		utils.CreateResponse(&response, http.StatusInternalServerError, map[string]interface{}{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	responseBody := map[string]interface{}{
		"error":   false,
		"message": "User promoted to service provider successfully",
	}

	utils.CreateResponse(&response, http.StatusOK, responseBody)
}
