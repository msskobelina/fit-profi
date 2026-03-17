package nutrition

import (
	"context"

	"github.com/msskobelina/fit-profi/internal/domain/model"
	"github.com/msskobelina/fit-profi/internal/domain/repository"
)

type UpdateEntryHandler interface {
	UpdateEntry(context.Context, UpdateEntryCommand) (*model.DiaryEntry, error)
}

type updateEntryService struct {
	repo repository.NutritionRepository
}

func NewUpdateEntryService(repo repository.NutritionRepository) UpdateEntryHandler {
	return &updateEntryService{repo: repo}
}

func (s *updateEntryService) UpdateEntry(ctx context.Context, cmd UpdateEntryCommand) (*model.DiaryEntry, error) {
	return s.repo.UpdateEntry(ctx, cmd.EntryID, cmd.UserID, model.DiaryEntry{
		MealType: cmd.MealType,
		Items:    cmd.Items,
	})
}
