package repository

import (
	"context"

	"github.com/msskobelina/fit-profi/internal/domain/model"
)

type ProfilesRepository interface {
	CreateUserProfile(ctx context.Context, p model.UserProfile) (*model.UserProfile, error)
	UpdateUserProfile(ctx context.Context, userID int, p model.UserProfile) (*model.UserProfile, error)
	GetUserProfileByUserID(ctx context.Context, userID int) (*model.UserProfile, error)
	CreateCoachProfile(ctx context.Context, p model.CoachProfile) (*model.CoachProfile, error)
	UpdateCoachProfile(ctx context.Context, userID int, p model.CoachProfile) (*model.CoachProfile, error)
	GetCoachProfileByUserID(ctx context.Context, userID int) (*model.CoachProfile, error)
}
