package interfaces

import (
	"context"

	"github.com/ESSantana/boilerplate-backend/internal/domain/models"
)

type CustomerRepository interface {
	GetCustomerLogin(ctx context.Context, email string, passwordHash string) (customer models.Customer, err error)
	GetCustomerById(ctx context.Context, id string) (user models.Customer, err error)
	GetCustomerByExternalID(ctx context.Context, externalID string) (customer models.Customer, err error)
	GetCustomerEmail(ctx context.Context, email string) (customer models.Customer, err error)
	GetAllCustomers(ctx context.Context) ([]models.Customer, error)
	CreateCustomer(ctx context.Context, customer models.Customer) error
	UpdateCustomer(ctx context.Context, customer models.Customer) error
	SoftDeleteCustomer(ctx context.Context, id string) error
}
