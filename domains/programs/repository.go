package programs

import (
	"context"

	"github.com/msskobelina/fit-profi/pkg/mysql"
	"gorm.io/gorm"
)

type Repository interface {
	CreateProgram(ctx context.Context, p *TrainingProgram) (*TrainingProgram, error)
	GetProgram(ctx context.Context, id int) (*TrainingProgram, error)
	ListProgramsByUser(ctx context.Context, userID int) ([]TrainingProgram, error)
	DeleteProgram(ctx context.Context, id int) error

	AddProgress(ctx context.Context, pr *ExerciseProgress) (*ExerciseProgress, error)
	GetProgramByExerciseID(ctx context.Context, exID int) (*TrainingProgram, error)
}

type gormRepo struct{ *mysql.MySQL }

func NewRepository(sql *mysql.MySQL) Repository {
	return &gormRepo{
		sql,
	}
}

func (r *gormRepo) CreateProgram(ctx context.Context, p *TrainingProgram) (*TrainingProgram, error) {
	if err := r.DB.WithContext(ctx).Create(p).Error; err != nil {
		return nil, err
	}

	return p, nil
}

func (r *gormRepo) GetProgram(ctx context.Context, id int) (*TrainingProgram, error) {
	var p TrainingProgram
	if err := r.DB.WithContext(ctx).
		Preload("Days.Exercises").
		Where("id = ?", id).
		First(&p).Error; err != nil {
		return nil, err
	}

	return &p, nil
}

func (r *gormRepo) ListProgramsByUser(ctx context.Context, userID int) ([]TrainingProgram, error) {
	var list []TrainingProgram
	if err := r.DB.WithContext(ctx).
		Where("user_id = ?", userID).
		Find(&list).Error; err != nil {
		return nil, err
	}

	return list, nil
}

func (r *gormRepo) DeleteProgram(ctx context.Context, id int) error {
	return r.DB.WithContext(ctx).Delete(&TrainingProgram{}, id).Error
}

func (r *gormRepo) AddProgress(ctx context.Context, pr *ExerciseProgress) (*ExerciseProgress, error) {
	if err := r.DB.WithContext(ctx).Create(pr).Error; err != nil {
		return nil, err
	}

	return pr, nil
}

func (r *gormRepo) GetProgramByExerciseID(ctx context.Context, exID int) (*TrainingProgram, error) {
	var ex ProgramExercise
	if err := r.DB.WithContext(ctx).First(&ex, exID).Error; err != nil {
		return nil, err
	}
	var day ProgramDay
	if err := r.DB.WithContext(ctx).First(&day, ex.ProgramDayID).Error; err != nil {
		return nil, err
	}
	var p TrainingProgram
	if err := r.DB.WithContext(ctx).First(&p, day.ProgramID).Error; err != nil {
		return nil, err
	}

	return &p, nil
}

var (
	ErrNotFound = gorm.ErrRecordNotFound
)
