package interfaces

import (
	"context"

	"github.com/application-ellas/ella-backend/internal/domain/models"
)

type CustomerService interface {
	GetCustomerLogin(ctx context.Context, email, passwordHash string) (customer models.Customer, err error)
	GetCustomerById(ctx context.Context, id string) (user models.Customer, err error)
	GetCustomerByExternalId(ctx context.Context, externalID string) (user models.Customer, err error)
	GetCustomerByEmail(ctx context.Context, email string) (customer models.Customer, err error)
	GetAllCustomers(ctx context.Context) (users []models.Customer, err error)
	CreateCustomer(ctx context.Context, customer models.Customer) (customerCreated models.Customer, err error)
	UpdateCustomer(ctx context.Context, customer models.Customer) error
	SoftDeleteCustomer(ctx context.Context, id string) error
}
