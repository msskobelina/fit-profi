package authorize

import (
	"context"

	"golang.org/x/crypto/bcrypt"

	"github.com/msskobelina/fit-profi/internal/domain/repository"
	"github.com/msskobelina/fit-profi/pkg/access"
)

type ResetPasswordHandler interface {
	ResetPassword(ctx context.Context, cmd ResetPasswordCommand) error
}

type resetPasswordService struct {
	repo       repository.UsersRepository
	hmacSecret string
}

func NewResetPasswordService(repo repository.UsersRepository, hmacSecret string) ResetPasswordHandler {
	return &resetPasswordService{repo: repo, hmacSecret: hmacSecret}
}

func (s *resetPasswordService) ResetPassword(ctx context.Context, cmd ResetPasswordCommand) error {
	if _, err := access.DecodeToken(cmd.Token, s.hmacSecret); err != nil {
		return err
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(cmd.Password), 14)
	if err != nil {
		return err
	}
	if err = s.repo.ResetPassword(ctx, cmd.Token, string(hash)); err != nil {
		return err
	}

	return s.repo.DeleteToken(ctx, cmd.Token)
}
