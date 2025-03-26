package interfaces

type RepositoryManager interface {
	NewUserRepository() UserRepository
}
