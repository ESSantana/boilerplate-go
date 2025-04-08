package interfaces

type RepositoryManager interface {
	DatabaseHealthCheck() error
	NewCustomerRepository() CustomerRepository
}
