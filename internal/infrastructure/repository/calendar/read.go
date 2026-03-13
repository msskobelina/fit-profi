package calendar

import (
	"context"
	"time"

	"github.com/msskobelina/fit-profi/internal/domain/model"
)

func (r *gormRepo) ListAvailability(ctx context.Context, coachID int, from, to time.Time) ([]model.CoachAvailability, error) {
	var res []model.CoachAvailability
	err := r.db.WithContext(ctx).
		Where("coach_id = ? AND start_time >= ? AND end_time <= ?", coachID, from.UTC(), to.UTC()).
		Order("start_time asc").
		Find(&res).Error
	return res, err
}
