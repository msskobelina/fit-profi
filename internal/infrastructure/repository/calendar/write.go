package calendar

import (
	"context"

	"gorm.io/gorm"

	"github.com/msskobelina/fit-profi/internal/domain/model"
	domainRepo "github.com/msskobelina/fit-profi/internal/domain/repository"
	"github.com/msskobelina/fit-profi/pkg/mysql"
)

type gormRepo struct {
	db *gorm.DB
}

func NewRepository(sql *mysql.MySQL) domainRepo.CalendarRepository {
	return &gormRepo{db: sql.DB}
}

func (r *gormRepo) AddAvailability(ctx context.Context, rows []model.CoachAvailability) error {
	return r.db.WithContext(ctx).Create(&rows).Error
}
