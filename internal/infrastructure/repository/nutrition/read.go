package nutrition

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"

	"github.com/msskobelina/fit-profi/internal/domain/model"
	utilsErrors "github.com/msskobelina/fit-profi/pkg/errors"
)

func (r *gormRepo) GetEntryByID(ctx context.Context, id, userID int) (*model.DiaryEntry, error) {
	var e model.DiaryEntry
	err := r.db.WithContext(ctx).
		Preload("Items").
		Where("id = ? AND user_id = ?", id, userID).
		First(&e).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &utilsErrors.Error{Message: "Entry not found"}
		}
		return nil, err
	}
	return &e, nil
}

func (r *gormRepo) ListEntriesByDate(ctx context.Context, userID int, date time.Time) ([]model.DiaryEntry, error) {
	var entries []model.DiaryEntry
	start := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	end := start.Add(24 * time.Hour)
	err := r.db.WithContext(ctx).
		Preload("Items").
		Where("user_id = ? AND date >= ? AND date < ?", userID, start, end).
		Find(&entries).Error
	return entries, err
}
