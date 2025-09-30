package calendar

import (
	"context"
	"time"

	"github.com/msskobelina/fit-profi/pkg/mysql"
)

type Repository interface {
	AddAvailability(ctx context.Context, rows []CoachAvailability) error
	ListAvailability(ctx context.Context, coachID int, from, to time.Time) ([]CoachAvailability, error)
}

type gormRepo struct{ *mysql.MySQL }

func NewRepository(db *mysql.MySQL) Repository { return &gormRepo{db} }

func (r *gormRepo) AddAvailability(ctx context.Context, rows []CoachAvailability) error {
	return r.DB.WithContext(ctx).Create(&rows).Error
}

func (r *gormRepo) ListAvailability(ctx context.Context, coachID int, from, to time.Time) ([]CoachAvailability, error) {
	var res []CoachAvailability
	err := r.DB.WithContext(ctx).
		Where("coach_id=? AND start_time>=? AND end_time<=?", coachID, from.UTC(), to.UTC()).
		Order("start_time asc").Find(&res).Error
	return res, err
}
