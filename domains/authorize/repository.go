package authorize

import (
	"context"
	"errors"
	"gorm.io/gorm"

	utilsErrors "github.com/msskobelina/fit-profi/pkg/errors"
	"github.com/msskobelina/fit-profi/pkg/mysql"
	"golang.org/x/crypto/bcrypt"
)

type Repository interface {
	CreateUser(ctx context.Context, inp *RegisterUserRequest) (*User, error)
	GetUser(ctx context.Context, email string) (*User, error)
	CreateToken(ctx context.Context, email, token string) error
	GetToken(ctx context.Context, token string) (*UserToken, error)
	DeleteToken(ctx context.Context, token string) error
	SaveRevokedToken(ctx context.Context, jti string, exp int64) error
	IsTokenRevoked(ctx context.Context, jti string) (bool, error)
	ResetPassword(ctx context.Context, token, hashed string) error
}

type gormRepo struct{ *mysql.MySQL }

func NewRepository(sql *mysql.MySQL) Repository {
	return &gormRepo{sql}
}

func (r *gormRepo) CreateUser(ctx context.Context, inp *RegisterUserRequest) (*User, error) {
	if u, _ := r.GetUser(ctx, inp.Email); u != nil {
		return nil, &utilsErrors.Error{Message: "User with this email already exists"}
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(inp.Password), 14)
	if err != nil {
		return nil, err
	}
	u := &User{FullName: inp.FullName, Email: inp.Email, Password: string(hash)}
	if err := r.DB.WithContext(ctx).Create(u).Error; err != nil {
		return nil, err
	}
	return u, nil
}

func (r *gormRepo) GetUser(ctx context.Context, email string) (*User, error) {
	var u User
	err := r.DB.WithContext(ctx).Where("email = ?", email).First(&u).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, &utilsErrors.Error{Message: "User not registered"}
		}
		return nil, err
	}
	return &u, nil
}

func (r *gormRepo) CreateToken(ctx context.Context, email, token string) error {
	return r.DB.WithContext(ctx).Create(&UserToken{Email: email, Token: token}).Error
}

func (r *gormRepo) GetToken(ctx context.Context, token string) (*UserToken, error) {
	var ut UserToken
	err := r.DB.WithContext(ctx).Where("token = ?", token).First(&ut).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &utilsErrors.Error{Message: "User with provided token not found"}
		}
		return nil, err
	}
	return &ut, nil
}

func (r *gormRepo) DeleteToken(ctx context.Context, token string) error {
	return r.DB.WithContext(ctx).Delete(&UserToken{}, "token = ?", token).Error
}

func (r *gormRepo) SaveRevokedToken(ctx context.Context, jti string, exp int64) error {
	return r.DB.WithContext(ctx).Create(&RevokedToken{JTI: jti, ExpiresAt: exp}).Error
}

func (r *gormRepo) IsTokenRevoked(ctx context.Context, jti string) (bool, error) {
	var n int64
	if err := r.DB.WithContext(ctx).Model(&RevokedToken{}).Where("jti = ?", jti).Count(&n).Error; err != nil {
		return false, err
	}
	return n > 0, nil
}

func (r *gormRepo) ResetPassword(ctx context.Context, token, hashed string) error {
	ut, err := r.GetToken(ctx, token)
	if err != nil {
		return err
	}
	return r.DB.WithContext(ctx).Model(&User{}).Where("email = ?", ut.Email).Update("password", hashed).Error
}
