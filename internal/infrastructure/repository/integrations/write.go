package integrations

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

func NewRepository(sql *mysql.MySQL) domainRepo.IntegrationsRepository {
	return &gormRepo{db: sql.DB}
}

func (r *gormRepo) Upsert(ctx context.Context, row model.UserIntegration) (*model.UserIntegration, error) {
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND provider = ?", row.UserID, row.Provider).
		Assign(&row).
		FirstOrCreate(&row).Error
	return &row, err
}
