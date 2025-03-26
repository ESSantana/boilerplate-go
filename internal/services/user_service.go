package services

import (
	"context"
	"strings"
	"time"

	"github.com/ESSantana/boilerplate-go/internal/domain/errors"
	"github.com/ESSantana/boilerplate-go/internal/domain/models"
	repo_interfaces "github.com/ESSantana/boilerplate-go/internal/repositories/interfaces"
	"github.com/ESSantana/boilerplate-go/internal/services/interfaces"
	"github.com/ESSantana/boilerplate-go/internal/utils"
	"github.com/ESSantana/boilerplate-go/packages/log"
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

	loc, err := time.LoadLocation("America/Sao_Paulo")
	if err != nil {
		svc.logger.Errorf("error at loading location: %s", err.Error())
		return user, errors.NewOperationError("error creating user registration")
	}
	now := time.Now().In(loc)

	user = models.User{
		ID:              id,
		Name:            name,
		Email:           email,
		ProviderOrigin:  provider,
		ExternalID:      externalID,
		ProfileImageURL: profileImageURL,
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	err = userRepo.CreateUser(ctx, user)
	if err != nil {
		svc.logger.Errorf("error at create user: %s", err.Error())
		return user, errors.NewOperationError("error creating user registration")
	}

	return user, nil
}
