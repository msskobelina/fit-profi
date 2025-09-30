package authorize

import (
	"context"
	"errors"
	"fmt"
	"github.com/msskobelina/fit-profi/api/emails"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/msskobelina/fit-profi/pkg/access"
	utilsErrors "github.com/msskobelina/fit-profi/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	Register(ctx context.Context, inp *RegisterUserRequest) (*AuthResponse, error)
	Login(ctx context.Context, inp *LoginUserRequest) (*AuthResponse, error)
	Logout(ctx context.Context, token string) error
	VerifyAccessToken(ctx context.Context, token string) (bool, int, string)
	SendEmail(ctx context.Context, email string) error
	ResetPassword(ctx context.Context, token, password string) error
}

type service struct {
	repo       Repository
	emails     emails.EmailsAPI
	hmacSecret string
	adminName  string
	adminEmail string
}

func NewService(repo Repository, emails emails.EmailsAPI, hmacSecret, adminName, adminEmail string) Service {
	return &service{
		repo:       repo,
		emails:     emails,
		hmacSecret: hmacSecret,
		adminName:  adminName,
		adminEmail: adminEmail,
	}
}

func (s *service) Register(ctx context.Context, inp *RegisterUserRequest) (*AuthResponse, error) {
	u, err := s.repo.CreateUser(ctx, inp)
	if err != nil {
		return nil, err
	}
	role := access.UserRoleUser
	if u.FullName == s.adminName && u.Email == s.adminEmail {
		role = access.UserRoleAdmin
	}
	now := time.Now()
	tok, err := access.EncodeToken(&access.Token{
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
	return &AuthResponse{
		Token:    tok,
		UserID:   u.ID,
		FullName: u.FullName,
		Email:    u.Email,
	}, nil
}

func (s *service) Login(ctx context.Context, inp *LoginUserRequest) (*AuthResponse, error) {
	u, err := s.repo.GetUser(ctx, inp.Email)
	if err != nil {
		return nil, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(inp.Password)); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return nil, &utilsErrors.Error{Message: "Wrong password"}
		}
		return nil, err
	}
	role := access.UserRoleUser
	if u.FullName == s.adminName && u.Email == s.adminEmail {
		role = access.UserRoleAdmin
	}
	now := time.Now()
	tok, err := access.EncodeToken(&access.Token{
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
	return &AuthResponse{
		Token:    tok,
		UserID:   u.ID,
		FullName: u.FullName,
		Email:    u.Email,
	}, nil
}

func (s *service) Logout(ctx context.Context, token string) error {
	t, err := access.DecodeToken(token, s.hmacSecret)
	if err != nil {
		return nil
	}
	return s.repo.SaveRevokedToken(ctx, t.ID, t.ExpiresAt.Unix())
}

func (s *service) VerifyAccessToken(ctx context.Context, token string) (bool, int, string) {
	t, err := access.DecodeToken(token, s.hmacSecret)
	if err != nil {
		return false, 0, ""
	}
	revoked, err := s.repo.IsTokenRevoked(ctx, t.ID)
	if err != nil || revoked {
		return false, 0, ""
	}
	return true, t.UserID, string(t.UserRole)
}

func (s *service) SendEmail(ctx context.Context, email string) error {
	u, err := s.repo.GetUser(ctx, email)
	if err != nil {
		return err
	}
	role := access.UserRoleUser
	if u.FullName == s.adminName && u.Email == s.adminEmail {
		role = access.UserRoleAdmin
	}
	now := time.Now()
	token, err := access.EncodeToken(&access.Token{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(15 * time.Minute)),
			NotBefore: jwt.NewNumericDate(now),
			IssuedAt:  jwt.NewNumericDate(now),
			Issuer:    email,
		},
		UserID:   u.ID,
		UserRole: role,
	}, s.hmacSecret)
	if err != nil {
		return err
	}
	if err := s.repo.CreateToken(ctx, email, token); err != nil {
		return err
	}
	return s.emails.SendEmail(ctx, emails.SendEmailInput{
		To:          email,
		Subject:     "FitProfi: reset password",
		ContentType: "text/html",
		Body:        fmt.Sprintf(`<h2>FitProfi: reset password</h2><p>Hello!</p><p>To reset your password, use this token:</p><p><b>%s</b></p>`, token),
	})
}

func (s *service) ResetPassword(ctx context.Context, token, password string) error {
	if _, err := access.DecodeToken(token, s.hmacSecret); err != nil {
		return err
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return err
	}
	if err := s.repo.ResetPassword(ctx, token, string(hash)); err != nil {
		return err
	}
	return s.repo.DeleteToken(ctx, token)
}
