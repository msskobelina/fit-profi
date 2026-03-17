package repository

import (
	"context"

	"github.com/msskobelina/fit-profi/internal/domain/model"
)

type UsersRepository interface {
	CreateUser(ctx context.Context, fullName, email, passwordHash string) (*model.User, error)
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
	CreateToken(ctx context.Context, email, token string) error
	GetToken(ctx context.Context, token string) (*model.UserToken, error)
	DeleteToken(ctx context.Context, token string) error
	SaveRevokedToken(ctx context.Context, jti string, exp int64) error
	IsTokenRevoked(ctx context.Context, jti string) (bool, error)
	ResetPassword(ctx context.Context, token, passwordHash string) error
}
