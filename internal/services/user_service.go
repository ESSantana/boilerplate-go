package services

import (
	"context"
	"strings"

	"github.com/application-ellas/ellas-backend/internal/domain/errors"
	"github.com/application-ellas/ellas-backend/internal/domain/models"
	repo_interfaces "github.com/application-ellas/ellas-backend/internal/repositories/interfaces"
	"github.com/application-ellas/ellas-backend/internal/services/interfaces"
	"github.com/application-ellas/ellas-backend/internal/utils"
	"github.com/application-ellas/ellas-backend/packages/log"
)

type userService struct {
	logger      log.Logger
	repoManager repo_interfaces.RepositoryManager
}

func newUserService(logger log.Logger, repoManager repo_interfaces.RepositoryManager) interfaces.UserService {
	return &userService{
		logger:      logger,
		repoManager: repoManager,
	}
}

func (svc *userService) GetUserByID(ctx context.Context, id string) (user models.User, err error) {
	userRepo := svc.repoManager.NewUserRepository()
	user, err = userRepo.GetUserById(ctx, id)
	if err != nil {
		return user, errors.NewNotFoundError("user not found")
	}
	return user, nil
}

func (svc *userService) GetUserByExternalID(ctx context.Context, externalID string) (user models.User, err error) {
	userRepo := svc.repoManager.NewUserRepository()
	user, err = userRepo.GetUserByExternalId(ctx, externalID)
	if err != nil {
		return user, errors.NewNotFoundError("user not found")
	}
	return user, nil
}

func (svc *userService) CreateUserIfNotExists(ctx context.Context, name, email, provider, externalID, profileImageURL string) (user models.User, err error) {
	id, err := utils.SHA1Hash(email)
	if err != nil {
		svc.logger.Errorf("error at SHA1Hash: %s", err.Error())
		return user, errors.NewOperationError("error creating user registration")
	}

	userRepo := svc.repoManager.NewUserRepository()
	user, err = userRepo.GetUserById(ctx, id)
	if err != nil && !strings.Contains(err.Error(), "error scanning user") {
		return user, errors.NewNotFoundError("user not found")
	}

	if user.ID != "" {
		return user, nil
	}

	user = models.User{
		ID:              id,
		Name:            name,
		Email:           email,
		ProviderOrigin:  provider,
		ExternalID:      externalID,
		ProfileImageURL: profileImageURL,
	}

	err = userRepo.CreateUser(ctx, user)
	if err != nil {
		svc.logger.Errorf("error at userRepo.CreateUser: %s", err.Error())
		return user, errors.NewOperationError("error creating user registration")
	}

	user, err = userRepo.GetUserById(ctx, id)
	if err != nil {
		svc.logger.Errorf("error at userRepo.GetUserById: %s", err.Error())
		return user, errors.NewOperationError("error getting user registration")
	}

	return user, nil
}
