package profiles

import (
	"context"

	"github.com/msskobelina/fit-profi/internal/domain/model"
	"github.com/msskobelina/fit-profi/internal/domain/repository"
)

type UpdateUserProfileHandler interface {
	UpdateUserProfile(context.Context, UpdateUserProfileCommand) (*model.UserProfile, error)
}

type updateUserProfileService struct {
	repo repository.ProfilesRepository
}

func NewUpdateUserProfileService(repo repository.ProfilesRepository) UpdateUserProfileHandler {
	return &updateUserProfileService{repo: repo}
}

func (s *updateUserProfileService) UpdateUserProfile(ctx context.Context, cmd UpdateUserProfileCommand) (*model.UserProfile, error) {
	return s.repo.UpdateUserProfile(ctx, cmd.UserID, model.UserProfile{
		FullName:    cmd.FullName,
		Age:         cmd.Age,
		WeightKg:    cmd.WeightKg,
		Goal:        model.Goal(cmd.Goal),
		Description: cmd.Description,
	})
}
