package interfaces

import (
	"context"

	"github.com/application-ellas/ellas-backend/internal/domain/models"
)

type CustomerService interface {
	GetCustomerById(ctx context.Context, id string) (user models.Customer, err error)
	GetCustomerByExternalId(ctx context.Context, externalID string) (user models.Customer, err error)
	GetAllCustomers(ctx context.Context) (users []models.Customer, err error)
	CreateCustomer(ctx context.Context, customer models.Customer) (customerCreated models.Customer, err error)
	UpdateCustomer(ctx context.Context, customer models.Customer) error
	SoftDeleteCustomer(ctx context.Context, id string) error
}
