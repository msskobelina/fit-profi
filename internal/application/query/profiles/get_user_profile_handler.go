package profiles

import (
	"context"

	"github.com/msskobelina/fit-profi/internal/domain/model"
	"github.com/msskobelina/fit-profi/internal/domain/repository"
)

type GetUserProfileHandler interface {
	GetUserProfile(context.Context, GetUserProfileQuery) (*model.UserProfile, error)
}

type getUserProfileService struct {
	repo repository.ProfilesRepository
}

func NewGetUserProfileService(repo repository.ProfilesRepository) GetUserProfileHandler {
	return &getUserProfileService{repo: repo}
}

func (s *getUserProfileService) GetUserProfile(
	ctx context.Context,
	q GetUserProfileQuery,
) (*model.UserProfile, error) {
	return s.repo.GetUserProfileByUserID(ctx, q.UserID)
}
