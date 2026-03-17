package nutrition

import (
	"context"

	"github.com/msskobelina/fit-profi/internal/domain/repository"
)

type DeleteEntryHandler interface {
	DeleteEntry(context.Context, DeleteEntryCommand) error
}

type deleteEntryService struct {
	repo repository.NutritionRepository
}

func NewDeleteEntryService(repo repository.NutritionRepository) DeleteEntryHandler {
	return &deleteEntryService{repo: repo}
}

func (s *deleteEntryService) DeleteEntry(ctx context.Context, cmd DeleteEntryCommand) error {
	return s.repo.DeleteEntry(ctx, cmd.EntryID, cmd.UserID)
}
