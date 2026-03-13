package profiles

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/msskobelina/fit-profi/internal/domain/model"
	domainRepo "github.com/msskobelina/fit-profi/internal/domain/repository"
	"github.com/msskobelina/fit-profi/pkg/mysql"
)

type gormRepo struct {
	db *gorm.DB
}

func NewRepository(sql *mysql.MySQL) domainRepo.ProfilesRepository {
	return &gormRepo{db: sql.DB}
}

func (r *gormRepo) CreateUserProfile(ctx context.Context, p model.UserProfile) (*model.UserProfile, error) {
	if err := r.db.WithContext(ctx).Create(&p).Error; err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *gormRepo) UpdateUserProfile(ctx context.Context, userID int, p model.UserProfile) (*model.UserProfile, error) {
	if err := r.db.WithContext(ctx).
		Model(&model.UserProfile{}).
		Where("user_id = ?", userID).
		Updates(&p).Error; err != nil {
		return nil, err
	}
	return r.GetUserProfileByUserID(ctx, userID)
}

func (r *gormRepo) CreateCoachProfile(ctx context.Context, p model.CoachProfile) (*model.CoachProfile, error) {
	if err := r.db.WithContext(ctx).Create(&p).Error; err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *gormRepo) UpdateCoachProfile(ctx context.Context, userID int, p model.CoachProfile) (*model.CoachProfile, error) {
	if err := r.db.WithContext(ctx).
		Session(&gorm.Session{FullSaveAssociations: true}).
		Where("user_id = ?", userID).
		Omit(clause.Associations).
		Updates(&p).Error; err != nil {
		return nil, err
	}
	return r.GetCoachProfileByUserID(ctx, userID)
}
