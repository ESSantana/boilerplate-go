package services

import (
	"context"
	"strings"

	"github.com/ESSantana/boilerplate-backend/internal/domain/constants"
	"github.com/ESSantana/boilerplate-backend/internal/domain/errors"
	"github.com/ESSantana/boilerplate-backend/internal/domain/models"
	repo_interfaces "github.com/ESSantana/boilerplate-backend/internal/repositories/interfaces"
	"github.com/ESSantana/boilerplate-backend/internal/services/interfaces"
	"github.com/ESSantana/boilerplate-backend/internal/utils"
	"github.com/ESSantana/boilerplate-backend/packages/log"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type customerService struct {
	logger      log.Logger
	repoManager repo_interfaces.RepositoryManager
	validate    *validator.Validate
}

func newCustomerService(logger log.Logger, repoManager repo_interfaces.RepositoryManager) interfaces.CustomerService {
	validate := validator.New(validator.WithRequiredStructEnabled())

	validate.RegisterValidation("required_if_match_format", utils.ValidateRequiredIfFieldMatchFormat)
	validate.RegisterValidation("required_if_not_match_format", utils.ValidateRequiredIfFieldNotMatchFormat)

	return &customerService{
		logger:      logger,
		repoManager: repoManager,
		validate:    validate,
	}
}

func (svc *customerService) GetCustomerLogin(ctx context.Context, email, passwordHash string) (customer models.Customer, err error) {
	customerRepo := svc.repoManager.NewCustomerRepository()

	if email == "" || passwordHash == "" {
		return customer, errors.NewValidationError("email and password are required")
	}

	customer, err = customerRepo.GetCustomerLogin(ctx, email, passwordHash)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return customer, errors.NewForbiddenError("email or password is incorrect")
		}
		return customer, err
	}
	return customer, nil
}

func (svc *customerService) GetCustomerById(ctx context.Context, id string) (customer models.Customer, err error) {
	customerRepo := svc.repoManager.NewCustomerRepository()
	customer, err = customerRepo.GetCustomerById(ctx, id)
	if err != nil {
		return customer, errors.NewNotFoundError("customer not found")
	}
	return customer, nil
}

func (svc *customerService) GetCustomerByExternalId(ctx context.Context, externalID string) (customer models.Customer, err error) {
	customerRepo := svc.repoManager.NewCustomerRepository()
	customer, err = customerRepo.GetCustomerByExternalID(ctx, externalID)

	if err != nil && !strings.Contains(err.Error(), "no rows in result set") {
		svc.logger.Errorf("error at customerRepo.GetCustomerByExternalID: %s", err.Error())
		return customer, err
	}
	return customer, nil
}

func (svc *customerService) GetCustomerByEmail(ctx context.Context, email string) (customer models.Customer, err error) {
	customerRepo := svc.repoManager.NewCustomerRepository()

	customerPersisted, err := customerRepo.GetCustomerEmail(ctx, email)
	if err != nil && !strings.Contains(err.Error(), "no rows in result set") {
		svc.logger.Errorf("error at customerRepo.GetCustomerEmail: %s", err.Error())
		return customer, err
	}

	return customerPersisted, nil
}

func (svc *customerService) GetAllCustomers(ctx context.Context) (customers []models.Customer, err error) {
	customerRepo := svc.repoManager.NewCustomerRepository()
	customers, err = customerRepo.GetAllCustomers(ctx)

	if err != nil && !strings.Contains(err.Error(), "no rows in result set") {
		svc.logger.Errorf("error at customerRepo.GetAllCustomers: %s", err.Error())
		return customers, err
	}
	return customers, nil
}

func (svc *customerService) CreateCustomer(ctx context.Context, customer models.Customer) (customerCreated models.Customer, err error) {
	if errValidation := svc.validate.Struct(customer); errValidation != nil {
		return customerCreated, utils.ParseValidatorErrorMessage(errValidation)
	}

	customerRepo := svc.repoManager.NewCustomerRepository()

	customer.ID = uuid.New().String()
	customer.ProviderOrigin = constants.ProviderOriginInternal

	customerPersisted, err := customerRepo.GetCustomerEmail(ctx, customer.Email)
	if err != nil && !strings.Contains(err.Error(), "no rows in result set") {
		svc.logger.Errorf("error at customerRepo.GetCustomerEmail: %s", err.Error())
		return customerCreated, err
	}

	if customerPersisted.ID != "" {
		return customerCreated, errors.NewValidationError("customer already exists")
	}

	err = customerRepo.CreateCustomer(ctx, customer)
	if err != nil {
		svc.logger.Errorf("error at customerRepo.CreateCustomer: %s", err.Error())
		return customerCreated, errors.NewOperationError("error creating customer registry")
	}

	return customer, nil
}

func (svc *customerService) UpdateCustomer(ctx context.Context, customer models.Customer) error {
	if errValidation := svc.validate.Struct(customer); errValidation != nil {
		return utils.ParseValidatorErrorMessage(errValidation)
	}

	customerRepo := svc.repoManager.NewCustomerRepository()

	_, err := customerRepo.GetCustomerById(ctx, customer.ID)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return errors.NewValidationError("customer not found")
		}
		return err
	}

	err = customerRepo.UpdateCustomer(ctx, customer)
	if err != nil {
		svc.logger.Errorf("error at customerRepo.UpdateCustomer: %s", err.Error())
		return errors.NewOperationError("error updating customer registry")
	}
	return nil
}

func (svc *customerService) SoftDeleteCustomer(ctx context.Context, id string) error {
	customerRepo := svc.repoManager.NewCustomerRepository()

	_, err := customerRepo.GetCustomerById(ctx, id)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return errors.NewValidationError("customer not found")
		}
		return err
	}

	err = customerRepo.SoftDeleteCustomer(ctx, id)
	if err != nil {
		svc.logger.Errorf("error at customerRepo.SoftDeleteCustomer: %s", err.Error())
		return errors.NewOperationError("error deleting customer registry")
	}
	return nil
}
