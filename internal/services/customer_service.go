package services

import (
	"context"
	"strings"

	"github.com/application-ellas/ellas-backend/internal/domain/constants"
	"github.com/application-ellas/ellas-backend/internal/domain/errors"
	"github.com/application-ellas/ellas-backend/internal/domain/models"
	repo_interfaces "github.com/application-ellas/ellas-backend/internal/repositories/interfaces"
	"github.com/application-ellas/ellas-backend/internal/services/interfaces"
	"github.com/application-ellas/ellas-backend/packages/log"
	"github.com/google/uuid"
)

type customerService struct {
	logger      log.Logger
	repoManager repo_interfaces.RepositoryManager
}

func newCustomerService(logger log.Logger, repoManager repo_interfaces.RepositoryManager) interfaces.CustomerService {
	return &customerService{
		logger:      logger,
		repoManager: repoManager,
	}
}

func (svc *customerService) GetCustomerById(ctx context.Context, id string) (user models.Customer, err error) {
	customerRepo := svc.repoManager.NewCustomerRepository()
	user, err = customerRepo.GetCustomerById(ctx, id)
	if err != nil {
		return user, errors.NewNotFoundError("customer not found")
	}
	return user, nil
}

func (svc *customerService) GetCustomerByExternalId(ctx context.Context, externalID string) (user models.Customer, err error) {
	customerRepo := svc.repoManager.NewCustomerRepository()
	user, err = customerRepo.GetCustomerByExternalID(ctx, externalID)
	if err != nil {
		return user, errors.NewNotFoundError("customer not found")
	}
	return user, nil
}

func (svc *customerService) GetAllCustomers(ctx context.Context) (users []models.Customer, err error) {
	customerRepo := svc.repoManager.NewCustomerRepository()
	users, err = customerRepo.GetAllCustomers(ctx)
	if err != nil {
		return users, errors.NewNotFoundError("any customer found")
	}
	return users, nil
}

func (svc *customerService) CreateCustomer(ctx context.Context, customer models.Customer) (customerCreated models.Customer, err error) {
	customerRepo := svc.repoManager.NewCustomerRepository()
	customerPersisted, err := customerRepo.GetCustomerEmail(ctx, customer.Email)
	if err != nil && !strings.Contains(err.Error(), "error scanning customer") {
		return customerCreated, errors.NewNotFoundError("customer not found")
	}

	if customerPersisted.ID != "" {
		return customerCreated, errors.NewValidationError("customer already exists")
	}

	customer.ID = uuid.New().String()
	customer.ExternalID = &customer.ID
	customer.ProviderOrigin = constants.ProviderOriginInternal

	if errValidation := customer.Validate(); errValidation != nil {
		return customerCreated, errors.NewValidationError(errValidation.Error())
	}

	err = customerRepo.CreateCustomer(ctx, customer)
	if err != nil {
		svc.logger.Errorf("error at customerRepo.CreateCustomer: %s", err.Error())
		return customerCreated, errors.NewOperationError("error creating customer registry")
	}

	return customer, nil
}

func (svc *customerService) UpdateCustomer(ctx context.Context, customer models.Customer) error {
	customerRepo := svc.repoManager.NewCustomerRepository()

	if errValidation := customer.Validate(); errValidation != nil {
		return errors.NewValidationError(errValidation.Error())
	}

	err := customerRepo.UpdateCustomer(ctx, customer)
	if err != nil {
		svc.logger.Errorf("error at ucustomerRepo.UpdateCustomer: %s", err.Error())
		return errors.NewOperationError("error updating customer registry")
	}
	return nil
}

func (svc *customerService) SoftDeleteCustomer(ctx context.Context, id string) error {
	customerRepo := svc.repoManager.NewCustomerRepository()
	err := customerRepo.SoftDeleteCustomer(ctx, id)
	if err != nil {
		svc.logger.Errorf("error at customerRepo.SoftDeleteCustomer: %s", err.Error())
		return errors.NewOperationError("error deleting customer registry")
	}
	return nil
}
