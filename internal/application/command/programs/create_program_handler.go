package programs

import (
	"context"

	"github.com/msskobelina/fit-profi/internal/domain/model"
	"github.com/msskobelina/fit-profi/internal/domain/repository"
)

type CreateProgramHandler interface {
	CreateProgram(context.Context, CreateProgramCommand) (*model.TrainingProgram, error)
}

type createProgramService struct {
	repo repository.ProgramsRepository
}

func NewCreateProgramService(repo repository.ProgramsRepository) CreateProgramHandler {
	return &createProgramService{repo: repo}
}

func (s *createProgramService) CreateProgram(ctx context.Context, cmd CreateProgramCommand) (*model.TrainingProgram, error) {
	return s.repo.CreateProgram(ctx, model.TrainingProgram{
		UserID:      cmd.UserID,
		Title:       cmd.Title,
		Description: cmd.Description,
		Days:        cmd.Days,
	})
}
