package repository

import (
	"context"
	"time"

	"github.com/msskobelina/fit-profi/internal/domain/model"
)

type CalendarRepository interface {
	AddAvailability(ctx context.Context, rows []model.CoachAvailability) error
	ListAvailability(ctx context.Context, coachID int, from, to time.Time) ([]model.CoachAvailability, error)
}
