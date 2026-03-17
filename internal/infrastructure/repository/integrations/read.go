package integrations

import (
	"context"

	"github.com/msskobelina/fit-profi/internal/domain/model"
)

func (r *gormRepo) GetByUserAndProvider(ctx context.Context, userID int, provider model.Provider) (*model.UserIntegration, error) {
	var u model.UserIntegration
	if err := r.db.WithContext(ctx).Where("user_id = ? AND provider = ?", userID, provider).First(&u).Error; err != nil {
		return nil, err
	}
	return &u, nil
}
