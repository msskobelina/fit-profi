package authorize

import (
	"context"

	"github.com/msskobelina/fit-profi/internal/domain/repository"
	"github.com/msskobelina/fit-profi/pkg/access"
)

type LogoutUserHandler interface {
	Logout(ctx context.Context, cmd LogoutUserCommand) error
}

type logoutUserService struct {
	repo       repository.UsersRepository
	hmacSecret string
}

func NewLogoutUserService(repo repository.UsersRepository, hmacSecret string) LogoutUserHandler {
	return &logoutUserService{repo: repo, hmacSecret: hmacSecret}
}

func (s *logoutUserService) Logout(ctx context.Context, cmd LogoutUserCommand) error {
	t, err := access.DecodeToken(cmd.Token, s.hmacSecret)
	if err != nil {
		return nil
	}

	return s.repo.SaveRevokedToken(ctx, t.ID, t.ExpiresAt.Unix())
}
