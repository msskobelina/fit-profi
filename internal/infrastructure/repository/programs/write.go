package programs

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

func NewRepository(sql *mysql.MySQL) domainRepo.ProgramsRepository {
	return &gormRepo{db: sql.DB}
}

func (r *gormRepo) CreateProgram(ctx context.Context, p model.TrainingProgram) (*model.TrainingProgram, error) {
	if err := r.db.WithContext(ctx).Create(&p).Error; err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *gormRepo) DeleteProgram(ctx context.Context, id, userID int) error {
	return r.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", id, userID).
		Delete(&model.TrainingProgram{}).Error
}

func (r *gormRepo) TrackProgress(ctx context.Context, prog model.ExerciseProgress) (*model.ExerciseProgress, error) {
	if err := r.db.WithContext(ctx).Create(&prog).Error; err != nil {
		return nil, err
	}
	return &prog, nil
}
