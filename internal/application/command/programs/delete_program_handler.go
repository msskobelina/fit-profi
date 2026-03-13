package programs

import (
	"context"

	"github.com/msskobelina/fit-profi/internal/domain/repository"
)

type DeleteProgramHandler interface {
	DeleteProgram(context.Context, DeleteProgramCommand) error
}

type deleteProgramService struct {
	repo repository.ProgramsRepository
}

func NewDeleteProgramService(repo repository.ProgramsRepository) DeleteProgramHandler {
	return &deleteProgramService{repo: repo}
}

func (s *deleteProgramService) DeleteProgram(ctx context.Context, cmd DeleteProgramCommand) error {
	return s.repo.DeleteProgram(ctx, cmd.ProgramID, cmd.UserID)
}
