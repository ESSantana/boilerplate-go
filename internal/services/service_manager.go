package services

import (
	repo_interfaces "github.com/ESSantana/boilerplate-backend/internal/repositories/interfaces"
	"github.com/ESSantana/boilerplate-backend/internal/services/interfaces"
	cache_interfaces "github.com/ESSantana/boilerplate-backend/packages/cache/interfaces"
)

type serviceManager struct {
	repoManager  repo_interfaces.RepositoryManager
	cacheManager cache_interfaces.CacheManager
}

func NewServiceManager(repoManager repo_interfaces.RepositoryManager, cacheManager cache_interfaces.CacheManager) interfaces.ServiceManager {
	return &serviceManager{
		repoManager:  repoManager,
		cacheManager: cacheManager,
	}
}

func (sm *serviceManager) HealthCheck() (dbHealthStatus, cacheHealthStatus bool) {
	dbHealthStatus = true
	cacheHealthStatus = true

	err := sm.repoManager.DatabaseHealthCheck()
	if err != nil {
		dbHealthStatus = false
	}

	err = sm.cacheManager.CacheHealthCheck()
	if err != nil {
		cacheHealthStatus = false
	}

	return
}

func (sm *serviceManager) NewCustomerService() interfaces.CustomerService {
	return newCustomerService(sm.repoManager)
}
