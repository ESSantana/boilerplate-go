package services

import (
	"context"
	"github.com/application-ellas/ellas-backend/internal/domain/models"
	repo_interfaces "github.com/application-ellas/ellas-backend/internal/repositories/interfaces"
	"github.com/application-ellas/ellas-backend/internal/services/interfaces"

	"github.com/application-ellas/ellas-backend/packages/log"
)

type productService struct {
	logger      log.Logger
	repoManager repo_interfaces.RepositoryManager
}

func newProductService(logger log.Logger, repoManager repo_interfaces.RepositoryManager) interfaces.ProductService {
	return &productService{
		logger:      logger,
		repoManager: repoManager,
	}
}

func (svc *productService) Create(ctx context.Context, product models.Product) (err error) {
	return nil
}

func (svc *productService) Update(ctx context.Context, product models.Product) (err error) {
	return nil
}

func (svc *productService) GetOne(ctx context.Context, productID string) (err error) {
	return nil
}

func (svc *productService) GetFiltered(ctx context.Context) (err error) {
	return nil
}

func (svc *productService) Delete(ctx context.Context, productID string) (err error) {
	return nil
}
