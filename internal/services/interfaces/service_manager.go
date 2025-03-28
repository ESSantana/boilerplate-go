package interfaces

type ServiceManager interface {
	NewUserService() UserService
	NewServiceProviderService() ServiceProviderService
}
