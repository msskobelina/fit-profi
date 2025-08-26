package repository

import (
	"context"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/msskobelina/fit-profi/internal/domain"
	"github.com/msskobelina/fit-profi/internal/service"
	"github.com/msskobelina/fit-profi/pkg/mysql"
)

type UsersRepo struct{ *mysql.MySQL }

func NewUsersRepo(mysql *mysql.MySQL) *UsersRepo {
	return &UsersRepo{mysql}
}

func (r *UsersRepo) CreateUser(ctx context.Context, inp *service.RegisterUserInput) (*domain.User, error) {
	if u, _ := r.GetUser(ctx, inp.Email); u != nil {
		return nil, &service.Error{Message: "User with this email already exists"}
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(inp.Password), 14)
	if err != nil {
		return nil, err
	}
	user := &domain.User{FullName: inp.FullName, Email: inp.Email, Password: string(hash)}
	if err := r.DB.WithContext(ctx).Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UsersRepo) GetUser(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User
	err := r.DB.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &service.Error{Message: "User not registered"}
		}
		return nil, err
	}
	return &user, nil
}

func (r *UsersRepo) CreateToken(ctx context.Context, inp service.GenerateTokenInput) error {
	return r.DB.WithContext(ctx).Create(&domain.UserToken{Email: inp.Email, Token: inp.Token}).Error
}

func (r *UsersRepo) GetToken(ctx context.Context, token string) (*domain.UserToken, error) {
	var user domain.UserToken
	err := r.DB.WithContext(ctx).
		Model(domain.UserToken{}).
		Where("token = ?", token).
		First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &service.Error{Message: "User with provided token not found"}
		}
		return nil, err
	}
	return &user, nil
}

func (r *UsersRepo) SaveRevokedToken(ctx context.Context, jti string, exp int64) error {
	return r.DB.WithContext(ctx).Create(&domain.RevokedToken{
		JTI: jti, ExpiresAt: exp,
	}).Error
}

func (r *UsersRepo) IsTokenRevoked(ctx context.Context, jti string) (bool, error) {
	var n int64
	err := r.DB.WithContext(ctx).Model(&domain.RevokedToken{}).
		Where("jti = ?", jti).Count(&n).Error
	return n > 0, err
}

func (r *UsersRepo) DeleteToken(ctx context.Context, token string) error {
	return r.DB.WithContext(ctx).Delete(&domain.UserToken{}, "token = ?", token).Error
}

func (r *UsersRepo) ResetPassword(ctx context.Context, inp *service.ResetPasswordInput) error {
	ut, err := r.GetToken(ctx, inp.Token)
	if err != nil {
		return err
	}
	return r.DB.WithContext(ctx).
		Model(&domain.User{}).
		Where("email = ?", ut.Email).
		Update("password", inp.Password).Error
}
