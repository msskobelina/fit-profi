package repository

import (
	"context"

	"github.com/msskobelina/fit-profi/internal/domain/model"
)

type IntegrationsRepository interface {
	Upsert(ctx context.Context, row model.UserIntegration) (*model.UserIntegration, error)
	GetByUserAndProvider(ctx context.Context, userID int, provider model.Provider) (*model.UserIntegration, error)
}
