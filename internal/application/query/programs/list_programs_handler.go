package programs

import (
	"context"

	"github.com/msskobelina/fit-profi/internal/domain/model"
	"github.com/msskobelina/fit-profi/internal/domain/repository"
)

type ListProgramsHandler interface {
	ListPrograms(context.Context, ListProgramsQuery) ([]model.TrainingProgram, error)
}

type listProgramsService struct {
	repo repository.ProgramsRepository
}

func NewListProgramsService(repo repository.ProgramsRepository) ListProgramsHandler {
	return &listProgramsService{repo: repo}
}

func (s *listProgramsService) ListPrograms(ctx context.Context, q ListProgramsQuery) ([]model.TrainingProgram, error) {
	return s.repo.ListProgramsByUserID(ctx, q.UserID)
}
