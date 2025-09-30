package integrations

import (
	"context"
	"github.com/msskobelina/fit-profi/pkg/mysql"
)

type Repository interface {
	Upsert(ctx context.Context, row *UserIntegration) (*UserIntegration, error)
	GetByUser(ctx context.Context, userID int, provider Provider) (*UserIntegration, error)
}

type gormRepo struct{ *mysql.MySQL }

func NewRepository(db *mysql.MySQL) Repository {
	return &gormRepo{db}
}

func (r *gormRepo) Upsert(ctx context.Context, row *UserIntegration) (*UserIntegration, error) {
	err := r.DB.WithContext(ctx).
		Where("user_id=? AND provider=?", row.UserID, row.Provider).
		Assign(row).
		FirstOrCreate(row).Error
	return row, err
}

func (r *gormRepo) GetByUser(ctx context.Context, userID int, provider Provider) (*UserIntegration, error) {
	var u UserIntegration
	if err := r.DB.WithContext(ctx).Where("user_id=? AND provider=?", userID, provider).First(&u).Error; err != nil {
		return nil, err
	}
	return &u, nil
}
