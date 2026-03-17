package programs

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/msskobelina/fit-profi/internal/domain/model"
	utilsErrors "github.com/msskobelina/fit-profi/pkg/errors"
)

func (r *gormRepo) GetProgramByID(ctx context.Context, id int) (*model.TrainingProgram, error) {
	var p model.TrainingProgram
	err := r.db.WithContext(ctx).
		Preload("Days.Exercises").
		First(&p, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &utilsErrors.Error{Message: "Program not found"}
		}
		return nil, err
	}
	return &p, nil
}

func (r *gormRepo) ListProgramsByUserID(ctx context.Context, userID int) ([]model.TrainingProgram, error) {
	var ps []model.TrainingProgram
	err := r.db.WithContext(ctx).
		Preload("Days.Exercises").
		Where("user_id = ?", userID).
		Find(&ps).Error
	return ps, err
}
