package interfaces

type RepositoryManager interface {
	DatabaseHealthCheck() error
	NewUserRepository() UserRepository
}
