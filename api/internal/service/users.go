package service

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/msskobelina/fit-profi/pkg/access"
	"golang.org/x/crypto/bcrypt"
)

type UsersService struct {
	repos Repositories
	apis  APIs
}

func NewUserService(repos Repositories, apis APIs) *UsersService {
	return &UsersService{repos: repos, apis: apis}
}

func (us *UsersService) RegisterUser(ctx context.Context, inp *RegisterUserInput) (RegisterUserOutput, error) {
	user, err := us.repos.Users.CreateUser(ctx, inp)
	if err != nil {
		return RegisterUserOutput{}, err
	}

	role := access.UserRoleUser
	if user.FullName == os.Getenv("ADMIN_USER_FULLNAME") && user.Email == os.Getenv("ADMIN_USER_EMAIL") {
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
		UserID:   user.ID,
		UserRole: role,
	}, os.Getenv("HMAC_SECRET"))
	if err != nil {
		return RegisterUserOutput{}, err
	}

	return RegisterUserOutput{
		Token:    token,
		UserID:   user.ID,
		FullName: user.FullName,
		Email:    user.Email,
	}, nil
}

func (us *UsersService) LoginUser(ctx context.Context, inp *LoginUserInput) (LoginUserOutput, error) {
	user, err := us.repos.Users.GetUser(ctx, inp.Email)
	if err != nil {
		return LoginUserOutput{}, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(inp.Password)); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return LoginUserOutput{}, &Error{Message: "Wrong password"}
		}
		return LoginUserOutput{}, err
	}

	role := access.UserRoleUser
	if user.FullName == os.Getenv("ADMIN_USER_FULLNAME") && user.Email == os.Getenv("ADMIN_USER_EMAIL") {
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
		UserID:   user.ID,
		UserRole: role,
	}, os.Getenv("HMAC_SECRET"))
	if err != nil {
		return LoginUserOutput{}, err
	}

	return LoginUserOutput{
		Token:    token,
		UserID:   user.ID,
		FullName: user.FullName,
		Email:    user.Email,
	}, nil
}

func (us *UsersService) Logout(ctx context.Context, token string) error {
	t, err := access.DecodeToken(token, os.Getenv("HMAC_SECRET"))
	if err != nil {
		return nil
	}

	exp := t.ExpiresAt.Unix()

	return us.repos.Users.SaveRevokedToken(ctx, t.ID, exp)
}

func (us *UsersService) VerifyAccessToken(ctx context.Context, token string) (bool, int, string) {
	t, err := access.DecodeToken(token, os.Getenv("HMAC_SECRET"))
	if err != nil {
		return false, 0, ""
	}

	revoked, err := us.repos.Users.IsTokenRevoked(ctx, t.ID)
	if err != nil || revoked {
		return false, 0, ""
	}

	return true, t.UserID, string(t.UserRole)
}

func (us *UsersService) SendEmail(ctx context.Context, inp *SendUserEmailInput) error {
	user, err := us.repos.Users.GetUser(ctx, inp.Email)
	if err != nil {
		return err
	}
	role := access.UserRoleUser
	if user.FullName == os.Getenv("ADMIN_USER_FULLNAME") && user.Email == os.Getenv("ADMIN_USER_EMAIL") {
		role = access.UserRoleAdmin
	}

	now := time.Now()
	token, err := access.EncodeToken(&access.Token{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(15 * time.Minute)),
			NotBefore: jwt.NewNumericDate(now),
			IssuedAt:  jwt.NewNumericDate(now),
			Issuer:    inp.Email,
		},
		UserID:   user.ID,
		UserRole: role,
	}, os.Getenv("HMAC_SECRET"))
	if err != nil {
		return err
	}

	if err := us.repos.Users.CreateToken(ctx, GenerateTokenInput{Email: inp.Email, Token: token}); err != nil {
		return err
	}

	return us.apis.Emails.SendEmail(ctx, SendEmailInput{
		To:          inp.Email,
		Subject:     "FitProfi: reset password",
		ContentType: "text/html",
		Body: fmt.Sprintf(`
			<h2>FitProfi: reset password</h2>
			<p>Hello!</p>
			<p>To reset your password, use this token:</p>
			<p><b>%s</b></p>
		`, token),
	})
}

func (us *UsersService) ResetPassword(ctx context.Context, inp *ResetPasswordInput) error {
	if _, err := us.repos.Users.GetToken(ctx, inp.Token); err != nil {
		return err
	}
	if _, err := access.DecodeToken(inp.Token, os.Getenv("HMAC_SECRET")); err != nil {
		return err
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(inp.Password), 14)
	if err != nil {
		return err
	}
	if err := us.repos.Users.ResetPassword(ctx, &ResetPasswordInput{
		Token:    inp.Token,
		Password: string(hash),
	}); err != nil {
		return err
	}

	return us.repos.Users.DeleteToken(ctx, inp.Token)
}
