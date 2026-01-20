package profiles

import (
	"context"
	"errors"
	"fmt"

	"github.com/msskobelina/fit-profi/pkg/analytics"
)

type Service interface {
	CreateUserProfile(ctx context.Context, uid int, inp *CreateUserProfileRequest) (*UserProfile, error)
	UpdateUserProfile(ctx context.Context, uid int, patch *UpdateUserProfileRequest) (*UserProfile, error)
	GetUserProfile(ctx context.Context, requesterID int, requesterRole string, userID int) (*UserProfile, error)

	CreateCoachProfile(ctx context.Context, adminID int, inp *CreateCoachProfileRequest) (*CoachProfile, error)
	UpdateCoachProfile(ctx context.Context, adminID int, patch *UpdateCoachProfileRequest) (*CoachProfile, error)
	GetCoachProfile(ctx context.Context, requesterRole string, adminID int) (*CoachProfile, error)
}

type service struct {
	repo      Repository
	analytics analytics.Client
}

func NewService(repo Repository, analytics analytics.Client) Service {
	return &service{
		repo:      repo,
		analytics: analytics,
	}
}

func (s *service) CreateUserProfile(ctx context.Context, userId int, inp *CreateUserProfileRequest) (*UserProfile, error) {
	if inp == nil {
		return nil, errors.New("empty body")
	}
	if inp.Age < 0 || inp.Age > 120 {
		return nil, errors.New("invalid age")
	}
	if inp.WeightKg <= 0 || inp.WeightKg > 500 {
		return nil, errors.New("invalid weight")
	}

	if analyticsErr := s.analytics.Track(ctx, "User Profile Created", fmt.Sprintf("%d", userId), map[string]any{
		"user_id":  userId,
		"source":   "backend",
		"endpoint": "/profiles/user/create",
	}); analyticsErr != nil {
		fmt.Println("User profile update analytics error: ", analyticsErr)
	}

	return s.repo.CreateUserProfile(ctx, &UserProfile{
		UserID:      userId,
		FullName:    inp.FullName,
		Age:         inp.Age,
		WeightKg:    inp.WeightKg,
		Goal:        inp.Goal,
		Description: inp.Description,
	})
}

func (s *service) UpdateUserProfile(ctx context.Context, userId int, patch *UpdateUserProfileRequest) (*UserProfile, error) {
	if patch == nil {
		return nil, errors.New("empty body")
	}

	if analyticsErr := s.analytics.Track(ctx, "User Profile Updated", fmt.Sprintf("%d", userId), map[string]any{
		"user_id":  userId,
		"source":   "backend",
		"endpoint": "/profiles/user/update",
	}); analyticsErr != nil {
		fmt.Println("User profile update analytics error: ", analyticsErr)
	}

	return s.repo.UpdateUserProfile(ctx, userId, patch)
}

func (s *service) GetUserProfile(ctx context.Context, requesterID int, requesterRole string, userID int) (*UserProfile, error) {
	if requesterRole != "admin" && requesterID != userID {
		return nil, errors.New("forbidden")
	}

	return s.repo.GetUserProfileByUserID(ctx, userID)
}

func (s *service) CreateCoachProfile(ctx context.Context, adminID int, inp *CreateCoachProfileRequest) (*CoachProfile, error) {
	if inp == nil {
		return nil, errors.New("empty body")
	}
	cp := &CoachProfile{
		UserID:   adminID,
		FullName: inp.FullName,
		Category: inp.Category,
		Info:     inp.Info,
	}
	for _, a := range inp.Achievements {
		cp.Achievements = append(cp.Achievements, CoachAchievement{
			StartPeriod:    a.StartPeriod,
			EndPeriod:      a.EndPeriod,
			Title:          a.Title,
			CertificateURL: a.CertificateURL,
		})
	}
	for _, e := range inp.Education {
		cp.Education = append(cp.Education, CoachEducation{
			StartPeriod: e.StartPeriod,
			EndPeriod:   e.EndPeriod,
			Place:       e.Place,
			Description: e.Description,
		})
	}
	if analyticsErr := s.analytics.Track(ctx, "Coach Profile Created", fmt.Sprintf("%d", adminID), map[string]any{
		"user_id":  adminID,
		"source":   "backend",
		"endpoint": "/profiles/coach/create",
	}); analyticsErr != nil {
		fmt.Println("Coach profile create analytics error: ", analyticsErr)
	}

	return s.repo.CreateCoachProfile(ctx, cp)
}

func (s *service) UpdateCoachProfile(ctx context.Context, adminID int, patch *UpdateCoachProfileRequest) (*CoachProfile, error) {
	if patch == nil {
		return nil, errors.New("empty body")
	}

	if analyticsErr := s.analytics.Track(ctx, "Coach Profile Updated", fmt.Sprintf("%d", adminID), map[string]any{
		"user_id":  adminID,
		"source":   "backend",
		"endpoint": "/profiles/coach/update",
	}); analyticsErr != nil {
		fmt.Println("Coach profile update analytics error: ", analyticsErr)
	}

	return s.repo.UpdateCoachProfile(ctx, adminID, patch)
}

func (s *service) GetCoachProfile(ctx context.Context, requesterRole string, adminID int) (*CoachProfile, error) {
	if requesterRole != "admin" {
		return nil, errors.New("forbidden")
	}

	return s.repo.GetCoachProfileByUserID(ctx, adminID)
}
