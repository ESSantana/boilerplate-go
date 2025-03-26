package services

import (
	repo_interfaces "github.com/ESSantana/boilerplate-go/internal/repositories/interfaces"
	"github.com/ESSantana/boilerplate-go/internal/services/interfaces"
	"github.com/ESSantana/boilerplate-go/packages/log"
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
