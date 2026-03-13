package profiles

import (
	"context"

	"github.com/msskobelina/fit-profi/internal/domain/model"
	"github.com/msskobelina/fit-profi/internal/domain/repository"
)

type GetCoachProfileHandler interface {
	GetCoachProfile(context.Context, GetCoachProfileQuery) (*model.CoachProfile, error)
}

type getCoachProfileService struct {
	repo repository.ProfilesRepository
}

func NewGetCoachProfileService(repo repository.ProfilesRepository) GetCoachProfileHandler {
	return &getCoachProfileService{repo: repo}
}

func (s *getCoachProfileService) GetCoachProfile(ctx context.Context, q GetCoachProfileQuery) (*model.CoachProfile, error) {
	return s.repo.GetCoachProfileByUserID(ctx, q.UserID)
}
