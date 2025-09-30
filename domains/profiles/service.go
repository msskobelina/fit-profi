package profiles

import (
	"context"
	"errors"
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
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{
		repo,
	}
}

func (s *service) CreateUserProfile(ctx context.Context, uid int, inp *CreateUserProfileRequest) (*UserProfile, error) {
	if inp == nil {
		return nil, errors.New("empty body")
	}
	if inp.Age < 0 || inp.Age > 120 {
		return nil, errors.New("invalid age")
	}
	if inp.WeightKg <= 0 || inp.WeightKg > 500 {
		return nil, errors.New("invalid weight")
	}

	return s.repo.CreateUserProfile(ctx, &UserProfile{
		UserID:      uid,
		FullName:    inp.FullName,
		Age:         inp.Age,
		WeightKg:    inp.WeightKg,
		Goal:        inp.Goal,
		Description: inp.Description,
	})
}

func (s *service) UpdateUserProfile(ctx context.Context, uid int, patch *UpdateUserProfileRequest) (*UserProfile, error) {
	if patch == nil {
		return nil, errors.New("empty body")
	}

	return s.repo.UpdateUserProfile(ctx, uid, patch)
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

	return s.repo.CreateCoachProfile(ctx, cp)
}

func (s *service) UpdateCoachProfile(ctx context.Context, adminID int, patch *UpdateCoachProfileRequest) (*CoachProfile, error) {
	if patch == nil {
		return nil, errors.New("empty body")
	}

	return s.repo.UpdateCoachProfile(ctx, adminID, patch)
}

func (s *service) GetCoachProfile(ctx context.Context, requesterRole string, adminID int) (*CoachProfile, error) {
	if requesterRole != "admin" {
		return nil, errors.New("forbidden")
	}

	return s.repo.GetCoachProfileByUserID(ctx, adminID)
}
