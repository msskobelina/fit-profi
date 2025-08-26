package service

import (
	"context"
	"github.com/msskobelina/fit-profi/internal/domain"
)

///
/// REPOSITORIES
///

type Repositories struct {
	Users UsersRepo
}

type UsersRepo interface {
	CreateUser(ctx context.Context, inp *RegisterUserInput) (*domain.User, error)
	GetUser(ctx context.Context, email string) (*domain.User, error)
	SaveRevokedToken(ctx context.Context, jti string, exp int64) error
	IsTokenRevoked(ctx context.Context, jti string) (bool, error)
	CreateToken(ctx context.Context, inp GenerateTokenInput) error
	GetToken(ctx context.Context, token string) (*domain.UserToken, error)
	DeleteToken(ctx context.Context, token string) error
	ResetPassword(ctx context.Context, inp *ResetPasswordInput) error
}

type GenerateTokenInput struct {
	Email string `json:"email"`
	Token string `json:"token"`
}
type ResetPasswordInput struct {
	Token    string
	Password string
}
