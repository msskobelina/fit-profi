package nutrition

import (
	"context"

	"github.com/msskobelina/fit-profi/internal/domain/model"
	"github.com/msskobelina/fit-profi/internal/domain/repository"
)

type GetEntryHandler interface {
	GetEntry(context.Context, GetEntryQuery) (*model.DiaryEntry, error)
}

type getEntryService struct {
	repo repository.NutritionRepository
}

func NewGetEntryService(repo repository.NutritionRepository) GetEntryHandler {
	return &getEntryService{repo: repo}
}

func (s *getEntryService) GetEntry(ctx context.Context, q GetEntryQuery) (*model.DiaryEntry, error) {
	return s.repo.GetEntryByID(ctx, q.EntryID, q.UserID)
}
