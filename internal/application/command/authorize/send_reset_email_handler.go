package authorize

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/msskobelina/fit-profi/internal/domain/repository"
	"github.com/msskobelina/fit-profi/internal/infrastructure/email"
	"github.com/msskobelina/fit-profi/pkg/access"
)

type SendResetEmailHandler interface {
	SendResetEmail(ctx context.Context, cmd SendResetEmailCommand) error
}

type sendResetEmailService struct {
	repo       repository.UsersRepository
	email      email.Sender
	hmacSecret string
	adminName  string
	adminEmail string
}

func NewSendResetEmailService(
	repo repository.UsersRepository,
	sender email.Sender,
	hmacSecret, adminName, adminEmail string,
) SendResetEmailHandler {
	return &sendResetEmailService{
		repo:       repo,
		email:      sender,
		hmacSecret: hmacSecret,
		adminName:  adminName,
		adminEmail: adminEmail,
	}
}

func (s *sendResetEmailService) SendResetEmail(ctx context.Context, cmd SendResetEmailCommand) error {
	u, err := s.repo.GetUserByEmail(ctx, cmd.Email)
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
			Issuer:    cmd.Email,
		},
		UserID:   u.ID,
		UserRole: role,
	}, s.hmacSecret)
	if err != nil {
		return err
	}

	if err = s.repo.CreateToken(ctx, cmd.Email, token); err != nil {
		return err
	}

	return s.email.Send(ctx, email.SendInput{
		To:          cmd.Email,
		Subject:     "FitProfi: reset password",
		ContentType: "text/html",
		Body:        fmt.Sprintf(`<h2>FitProfi: reset password</h2><p>Hello!</p><p>To reset your password, use this token:</p><p><b>%s</b></p>`, token),
	})
}
