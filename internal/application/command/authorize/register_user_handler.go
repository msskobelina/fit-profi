package authorize

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/msskobelina/fit-profi/internal/domain/repository"
	"github.com/msskobelina/fit-profi/pkg/access"
	"github.com/msskobelina/fit-profi/pkg/analytics"
	"github.com/msskobelina/fit-profi/pkg/metric"
)

type RegisterUserHandler interface {
	Register(ctx context.Context, cmd RegisterUserCommand) (*RegisterUserResult, error)
}

type registerUserService struct {
	repo       repository.UsersRepository
	analytics  analytics.Client
	metrics    *metric.Service
	hmacSecret string
	adminName  string
	adminEmail string
}

func NewRegisterUserService(
	repo repository.UsersRepository,
	analytics analytics.Client,
	metrics *metric.Service,
	hmacSecret, adminName, adminEmail string,
) RegisterUserHandler {
	return &registerUserService{
		repo:       repo,
		analytics:  analytics,
		metrics:    metrics,
		hmacSecret: hmacSecret,
		adminName:  adminName,
		adminEmail: adminEmail,
	}
}

func (s *registerUserService) Register(ctx context.Context, cmd RegisterUserCommand) (*RegisterUserResult, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(cmd.Password), 14)
	if err != nil {
		return nil, err
	}

	u, err := s.repo.CreateUser(ctx, cmd.FullName, cmd.Email, string(hash))
	if err != nil {
		return nil, err
	}

	s.metrics.TrackUserCreated("api")

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

	if analyticsErr := s.analytics.Track(ctx, "User Registered", fmt.Sprintf("%d", u.ID), map[string]any{
		"user_id": u.ID, "email": u.Email, "role": string(role),
	}); analyticsErr != nil {
		fmt.Println("Register analytics error:", analyticsErr)
	}

	return &RegisterUserResult{
		Token:    token,
		UserID:   u.ID,
		FullName: u.FullName,
		Email:    u.Email,
	}, nil
}
