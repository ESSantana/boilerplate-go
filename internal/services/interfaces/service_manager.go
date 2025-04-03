package interfaces

type ServiceManager interface {
	HealthCheck() (dbHealthStatus, cacheHealthStatus bool)
	NewUserService() UserService
	NewServiceProviderService() ServiceProviderService
	NewProductService() ProductService
}
