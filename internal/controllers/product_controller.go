package controllers

import (
	"context"
	"encoding/json"
	"io"

	"net/http"

	"github.com/application-ellas/ellas-backend/internal/domain/constants"
	"github.com/application-ellas/ellas-backend/internal/domain/models"
	svc_interfaces "github.com/application-ellas/ellas-backend/internal/services/interfaces"
	"github.com/application-ellas/ellas-backend/internal/utils"
	cache_interfaces "github.com/application-ellas/ellas-backend/packages/cache/interfaces"
	"github.com/application-ellas/ellas-backend/packages/log"
)

type ProductController struct {
	logger         log.Logger
	serviceManager svc_interfaces.ServiceManager
	cacheManager   cache_interfaces.CacheManager
}

func NewServiceController(logger log.Logger, serviceManager svc_interfaces.ServiceManager, cacheManager cache_interfaces.CacheManager) ProductController {
	return ProductController{
		logger:         logger,
		serviceManager: serviceManager,
		cacheManager:   cacheManager,
	}
}

func (ctlr *ProductController) Create(response http.ResponseWriter, request *http.Request) {
	context, cancel := context.WithTimeout(request.Context(), constants.DefaultTimeout)
	defer cancel()

	body, err := io.ReadAll(request.Body)
	if err != nil {
		utils.CreateResponse(&response, http.StatusBadRequest, err)
		return
	}

	var product models.Product
	err = json.Unmarshal(body, &product)
	if err != nil {
		utils.CreateResponse(&response, http.StatusBadRequest, err)
		return
	}

	err = product.Validate(constants.ValidationTypeCreate)
	if err != nil {
		utils.CreateResponse(&response, http.StatusBadRequest, err)
		return
	}

	ctlr.logger.Debugf("product received: %v", product)

	productService := ctlr.serviceManager.NewProductService()
	err = productService.Create(context, product)
	if err != nil {
		utils.CreateResponse(&response, http.StatusInternalServerError, err)
		return
	}
	utils.CreateResponse(&response, http.StatusOK, nil, "Product created successfully")
}

func (ctlr *ProductController) Update(response http.ResponseWriter, request *http.Request) {
	context, cancel := context.WithTimeout(request.Context(), constants.DefaultTimeout)
	defer cancel()

	body, err := io.ReadAll(request.Body)
	if err != nil {
		utils.CreateResponse(&response, http.StatusBadRequest, err)
		return
	}

	var product models.Product
	err = json.Unmarshal(body, &product)
	if err != nil {
		utils.CreateResponse(&response, http.StatusBadRequest, err)
		return
	}

	err = product.Validate(constants.ValidationTypeUpdate)
	if err != nil {
		utils.CreateResponse(&response, http.StatusBadRequest, err)
		return
	}

	ctlr.logger.Debugf("product received: %v", product)

	productService := ctlr.serviceManager.NewProductService()
	err = productService.Update(context, product)
	if err != nil {
		utils.CreateResponse(&response, http.StatusInternalServerError, err)
		return
	}
	utils.CreateResponse(&response, http.StatusOK, nil, "Product updated successfully")
}
