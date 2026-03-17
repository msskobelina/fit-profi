package nutrition

import (
	"context"

	"gorm.io/gorm"

	"github.com/msskobelina/fit-profi/internal/domain/model"
	domainRepo "github.com/msskobelina/fit-profi/internal/domain/repository"
	"github.com/msskobelina/fit-profi/pkg/mysql"
)

type gormRepo struct {
	db *gorm.DB
}

func NewRepository(sql *mysql.MySQL) domainRepo.NutritionRepository {
	return &gormRepo{db: sql.DB}
}

func (r *gormRepo) CreateEntry(ctx context.Context, e model.DiaryEntry) (*model.DiaryEntry, error) {
	if err := r.db.WithContext(ctx).Create(&e).Error; err != nil {
		return nil, err
	}
	return &e, nil
}

func (r *gormRepo) UpdateEntry(ctx context.Context, id, userID int, e model.DiaryEntry) (*model.DiaryEntry, error) {
	if err := r.db.WithContext(ctx).
		Model(&model.DiaryEntry{}).
		Where("id = ? AND user_id = ?", id, userID).
		Updates(&e).Error; err != nil {
		return nil, err
	}
	return r.GetEntryByID(ctx, id, userID)
}

func (r *gormRepo) DeleteEntry(ctx context.Context, id, userID int) error {
	return r.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", id, userID).
		Delete(&model.DiaryEntry{}).Error
}
