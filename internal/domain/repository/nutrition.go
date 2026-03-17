package repository

import (
	"context"
	"time"

	"github.com/msskobelina/fit-profi/internal/domain/model"
)

type NutritionRepository interface {
	CreateEntry(ctx context.Context, e model.DiaryEntry) (*model.DiaryEntry, error)
	GetEntryByID(ctx context.Context, id, userID int) (*model.DiaryEntry, error)
	ListEntriesByDate(ctx context.Context, userID int, date time.Time) ([]model.DiaryEntry, error)
	UpdateEntry(ctx context.Context, id, userID int, e model.DiaryEntry) (*model.DiaryEntry, error)
	DeleteEntry(ctx context.Context, id, userID int) error
}
