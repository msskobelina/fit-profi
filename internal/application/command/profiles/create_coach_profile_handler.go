package profiles

import (
	"context"

	"github.com/msskobelina/fit-profi/internal/domain/model"
	"github.com/msskobelina/fit-profi/internal/domain/repository"
)

type CreateCoachProfileHandler interface {
	CreateCoachProfile(context.Context, CreateCoachProfileCommand) (*model.CoachProfile, error)
}

type createCoachProfileService struct {
	repo repository.ProfilesRepository
}

func NewCreateCoachProfileService(repo repository.ProfilesRepository) CreateCoachProfileHandler {
	return &createCoachProfileService{repo: repo}
}

func (s *createCoachProfileService) CreateCoachProfile(ctx context.Context, cmd CreateCoachProfileCommand) (*model.CoachProfile, error) {
	return s.repo.CreateCoachProfile(ctx, model.CoachProfile{
		UserID:       cmd.UserID,
		FullName:     cmd.FullName,
		Category:     model.CoachCategory(cmd.Category),
		Info:         cmd.Info,
		Achievements: cmd.Achievements,
		Education:    cmd.Education,
	})
}
