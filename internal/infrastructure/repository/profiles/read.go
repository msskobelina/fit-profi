package profiles

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/msskobelina/fit-profi/internal/domain/model"
	utilsErrors "github.com/msskobelina/fit-profi/pkg/errors"
)

func (r *gormRepo) GetUserProfileByUserID(ctx context.Context, userID int) (*model.UserProfile, error) {
	var p model.UserProfile
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&p).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &utilsErrors.Error{Message: "Profile not found"}
		}
		return nil, err
	}
	return &p, nil
}

func (r *gormRepo) GetCoachProfileByUserID(ctx context.Context, userID int) (*model.CoachProfile, error) {
	var p model.CoachProfile
	err := r.db.WithContext(ctx).
		Preload("Achievements").
		Preload("Education").
		Where("user_id = ?", userID).First(&p).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &utilsErrors.Error{Message: "Coach profile not found"}
		}
		return nil, err
	}
	return &p, nil
}
