package interfaces

type ServiceManager interface {
	HealthCheck() (dbHealthStatus, cacheHealthStatus bool)
	NewCustomerService() CustomerService
}
