package authorize

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/msskobelina/fit-profi/internal/domain/repository"
	"github.com/msskobelina/fit-profi/pkg/access"
	"github.com/msskobelina/fit-profi/pkg/analytics"
	utilsErrors "github.com/msskobelina/fit-profi/pkg/errors"
	"github.com/msskobelina/fit-profi/pkg/metric"
)

type LoginUserHandler interface {
	Login(ctx context.Context, cmd LoginUserCommand) (*LoginUserResult, error)
}

type loginUserService struct {
	repo       repository.UsersRepository
	analytics  analytics.Client
	metrics    *metric.Service
	hmacSecret string
	adminName  string
	adminEmail string
}

func NewLoginUserService(
	repo repository.UsersRepository,
	analytics analytics.Client,
	metrics *metric.Service,
	hmacSecret, adminName, adminEmail string,
) LoginUserHandler {
	return &loginUserService{
		repo:       repo,
		analytics:  analytics,
		metrics:    metrics,
		hmacSecret: hmacSecret,
		adminName:  adminName,
		adminEmail: adminEmail,
	}
}

func (s *loginUserService) Login(ctx context.Context, cmd LoginUserCommand) (*LoginUserResult, error) {
	u, err := s.repo.GetUserByEmail(ctx, cmd.Email)
	if err != nil {
		s.metrics.TrackLoginFailed("user_not_found")
		return nil, err
	}

	if err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(cmd.Password)); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			s.metrics.TrackLoginFailed("wrong_password")
			return nil, &utilsErrors.Error{Message: "Wrong password"}
		}
		return nil, err
	}

	role := access.UserRoleUser
	if u.FullName == s.adminName && u.Email == s.adminEmail {
		role = access.UserRoleAdmin
	}

	now := time.Now()
	token, err := access.EncodeToken(&access.Token{
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        uuid.NewString(),
			ExpiresAt: jwt.NewNumericDate(now.Add(14 * 24 * time.Hour)),
			NotBefore: jwt.NewNumericDate(now),
			IssuedAt:  jwt.NewNumericDate(now),
			Issuer:    "fit-profi-api",
		},
		UserID:   u.ID,
		UserRole: role,
	}, s.hmacSecret)
	if err != nil {
		return nil, err
	}

	if analyticsErr := s.analytics.Track(ctx, "User Login", fmt.Sprintf("%d", u.ID), map[string]any{
		"user_id": u.ID, "email": u.Email, "role": string(role),
	}); analyticsErr != nil {
		fmt.Println("Login analytics error:", analyticsErr)
	}

	return &LoginUserResult{
		Token:    token,
		UserID:   u.ID,
		FullName: u.FullName,
		Email:    u.Email,
	}, nil
}
