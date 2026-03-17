package programs

import (
	"context"

	"github.com/msskobelina/fit-profi/internal/domain/model"
	"github.com/msskobelina/fit-profi/internal/domain/repository"
)

type GetProgramHandler interface {
	GetProgram(context.Context, GetProgramQuery) (*model.TrainingProgram, error)
}

type getProgramService struct {
	repo repository.ProgramsRepository
}

func NewGetProgramService(repo repository.ProgramsRepository) GetProgramHandler {
	return &getProgramService{repo: repo}
}

func (s *getProgramService) GetProgram(ctx context.Context, q GetProgramQuery) (*model.TrainingProgram, error) {
	return s.repo.GetProgramByID(ctx, q.ProgramID)
}
