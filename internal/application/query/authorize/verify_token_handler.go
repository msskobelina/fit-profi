package authorize

import (
	"context"

	"github.com/msskobelina/fit-profi/internal/domain/repository"
	"github.com/msskobelina/fit-profi/pkg/access"
)

type VerifyTokenHandler interface {
	VerifyToken(ctx context.Context, q VerifyTokenQuery) (*VerifyTokenResult, error)
}

type verifyTokenService struct {
	repo       repository.UsersRepository
	hmacSecret string
}

func NewVerifyTokenService(repo repository.UsersRepository, hmacSecret string) VerifyTokenHandler {
	return &verifyTokenService{repo: repo, hmacSecret: hmacSecret}
}

func (s *verifyTokenService) VerifyToken(ctx context.Context, q VerifyTokenQuery) (*VerifyTokenResult, error) {
	t, err := access.DecodeToken(q.Token, s.hmacSecret)
	if err != nil {
		return nil, err
	}
	revoked, err := s.repo.IsTokenRevoked(ctx, t.ID)
	if err != nil {
		return nil, err
	}
	if revoked {
		return nil, nil
	}

	return &VerifyTokenResult{
		UserID: t.UserID,
		Role:   string(t.UserRole),
	}, nil
}
