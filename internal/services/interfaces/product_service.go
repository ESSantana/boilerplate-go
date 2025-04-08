package interfaces

import (
	"context"

	"github.com/application-ellas/ella-backend/internal/domain/models"
)

type ProductService interface {
	Create(ctx context.Context, product models.Product) (err error)
	Update(ctx context.Context, product models.Product) (err error)
	GetOne(ctx context.Context, productID string) (err error)
	GetFiltered(ctx context.Context) (err error)
	Delete(ctx context.Context, productID string) (err error)
}
