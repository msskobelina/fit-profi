package profiles

import (
	"context"

	"github.com/msskobelina/fit-profi/internal/domain/model"
	"github.com/msskobelina/fit-profi/internal/domain/repository"
)

type UpdateCoachProfileHandler interface {
	UpdateCoachProfile(context.Context, UpdateCoachProfileCommand) (*model.CoachProfile, error)
}

type updateCoachProfileService struct {
	repo repository.ProfilesRepository
}

func NewUpdateCoachProfileService(repo repository.ProfilesRepository) UpdateCoachProfileHandler {
	return &updateCoachProfileService{repo: repo}
}

func (s *updateCoachProfileService) UpdateCoachProfile(ctx context.Context, cmd UpdateCoachProfileCommand) (*model.CoachProfile, error) {
	return s.repo.UpdateCoachProfile(ctx, cmd.UserID, model.CoachProfile{
		FullName:     cmd.FullName,
		Category:     model.CoachCategory(cmd.Category),
		Info:         cmd.Info,
		Achievements: cmd.Achievements,
		Education:    cmd.Education,
	})
}
