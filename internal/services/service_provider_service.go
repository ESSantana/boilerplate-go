package services

import (
	"context"
	"strings"

	"github.com/application-ellas/ellas-backend/internal/domain/errors"
	"github.com/application-ellas/ellas-backend/internal/domain/models"
	repo_interfaces "github.com/application-ellas/ellas-backend/internal/repositories/interfaces"
	"github.com/application-ellas/ellas-backend/internal/services/interfaces"

	"github.com/application-ellas/ellas-backend/packages/log"
)

type serviceProviderService struct {
	logger      log.Logger
	repoManager repo_interfaces.RepositoryManager
}

func newServiceProviderService(logger log.Logger, repoManager repo_interfaces.RepositoryManager) interfaces.ServiceProviderService {
	return &serviceProviderService{
		logger:      logger,
		repoManager: repoManager,
	}
}

func (svc *serviceProviderService) PromoteUserToServiceProvider(ctx context.Context, userID string) (err error) {
	userRepo := svc.repoManager.NewUserRepository()
	userRole, err := userRepo.GetUserRoleByUserIdAndRole(ctx, userID, models.RoleServiceProvider)
	if err != nil && !strings.Contains(err.Error(), "error scanning user") {
		return errors.NewOperationError("error getting user role")
	}

	if userRole.ID != "" {
		return errors.NewValidationError("user is already a service provider")
	}

	err = userRepo.AppendRoleToUser(ctx, userID, models.RoleServiceProvider)
	if err != nil {
		svc.logger.Errorf("error at userRepo.AppendRoleToUser: %s", err.Error())
		return errors.NewOperationError("error creating user registration")
	}
	return nil
}
