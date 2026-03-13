package nutrition

import (
	"context"

	"github.com/msskobelina/fit-profi/internal/domain/model"
	"github.com/msskobelina/fit-profi/internal/domain/repository"
)

type ListEntriesHandler interface {
	ListEntries(context.Context, ListEntriesQuery) ([]model.DiaryEntry, error)
}

type listEntriesService struct {
	repo repository.NutritionRepository
}

func NewListEntriesService(repo repository.NutritionRepository) ListEntriesHandler {
	return &listEntriesService{repo: repo}
}

func (s *listEntriesService) ListEntries(ctx context.Context, q ListEntriesQuery) ([]model.DiaryEntry, error) {
	return s.repo.ListEntriesByDate(ctx, q.UserID, q.Date)
}
