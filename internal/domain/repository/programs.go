package repository

import (
	"context"

	"github.com/msskobelina/fit-profi/internal/domain/model"
)

type ProgramsRepository interface {
	CreateProgram(ctx context.Context, p model.TrainingProgram) (*model.TrainingProgram, error)
	GetProgramByID(ctx context.Context, id int) (*model.TrainingProgram, error)
	ListProgramsByUserID(ctx context.Context, userID int) ([]model.TrainingProgram, error)
	DeleteProgram(ctx context.Context, id, userID int) error
	TrackProgress(ctx context.Context, prog model.ExerciseProgress) (*model.ExerciseProgress, error)
}
