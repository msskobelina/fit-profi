package profiles

import (
	"context"

	"github.com/msskobelina/fit-profi/internal/domain/model"
	"github.com/msskobelina/fit-profi/internal/domain/repository"
)

type CreateUserProfileHandler interface {
	CreateUserProfile(context.Context, CreateUserProfileCommand) (*model.UserProfile, error)
}

type createUserProfileService struct {
	repo repository.ProfilesRepository
}

func NewCreateUserProfileService(repo repository.ProfilesRepository) CreateUserProfileHandler {
	return &createUserProfileService{repo: repo}
}

func (s *createUserProfileService) CreateUserProfile(
	ctx context.Context,
	cmd CreateUserProfileCommand,
) (*model.UserProfile, error) {

	profile := model.UserProfile{
		UserID:      cmd.UserID,
		FullName:    cmd.FullName,
		Age:         cmd.Age,
		WeightKg:    cmd.WeightKg,
		Goal:        model.Goal(cmd.Goal),
		Description: cmd.Description,
	}

	return s.repo.CreateUserProfile(ctx, profile)
}
