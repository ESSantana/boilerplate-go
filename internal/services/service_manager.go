package services

import (
	repo_interfaces "github.com/application-ellas/ellas-backend/internal/repositories/interfaces"
	"github.com/application-ellas/ellas-backend/internal/services/interfaces"
	"github.com/application-ellas/ellas-backend/packages/log"
)

type serviceManager struct {
	logger      log.Logger
	repoManager repo_interfaces.RepositoryManager
}

func NewServiceManager(logger log.Logger, repoManager repo_interfaces.RepositoryManager) interfaces.ServiceManager {
	return &serviceManager{
		logger:      logger,
		repoManager: repoManager,
	}
}

func (sm *serviceManager) NewUserService() interfaces.UserService {
	return newUserService(sm.logger, sm.repoManager)
}

func (sm *serviceManager) NewServiceProviderService() interfaces.ServiceProviderService {
	return newServiceProviderService(sm.logger, sm.repoManager)
}
