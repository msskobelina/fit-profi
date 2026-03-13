package nutrition

import (
	"context"

	"github.com/msskobelina/fit-profi/internal/domain/model"
	"github.com/msskobelina/fit-profi/internal/domain/repository"
)

type CreateEntryHandler interface {
	CreateEntry(context.Context, CreateEntryCommand) (*model.DiaryEntry, error)
}

type createEntryService struct {
	repo repository.NutritionRepository
}

func NewCreateEntryService(repo repository.NutritionRepository) CreateEntryHandler {
	return &createEntryService{repo: repo}
}

func (s *createEntryService) CreateEntry(ctx context.Context, cmd CreateEntryCommand) (*model.DiaryEntry, error) {
	return s.repo.CreateEntry(ctx, model.DiaryEntry{
		UserID:   cmd.UserID,
		Date:     cmd.Date,
		MealType: cmd.MealType,
		Items:    cmd.Items,
	})
}
