package authorize

import (
	"context"

	"gorm.io/gorm"

	"github.com/msskobelina/fit-profi/internal/domain/model"
	domainRepo "github.com/msskobelina/fit-profi/internal/domain/repository"
	utilsErrors "github.com/msskobelina/fit-profi/pkg/errors"
	"github.com/msskobelina/fit-profi/pkg/mysql"
)

type gormRepo struct {
	db *gorm.DB
}

func NewRepository(sql *mysql.MySQL) domainRepo.UsersRepository {
	return &gormRepo{db: sql.DB}
}

func (r *gormRepo) CreateUser(ctx context.Context, fullName, email, passwordHash string) (*model.User, error) {
	if u, _ := r.GetUserByEmail(ctx, email); u != nil {
		return nil, &utilsErrors.Error{Message: "User with this email already exists"}
	}
	u := &model.User{FullName: fullName, Email: email, Password: passwordHash}
	if err := r.db.WithContext(ctx).Create(u).Error; err != nil {
		return nil, err
	}
	return u, nil
}

func (r *gormRepo) CreateToken(ctx context.Context, email, token string) error {
	return r.db.WithContext(ctx).Create(&model.UserToken{Email: email, Token: token}).Error
}

func (r *gormRepo) DeleteToken(ctx context.Context, token string) error {
	return r.db.WithContext(ctx).Delete(&model.UserToken{}, "token = ?", token).Error
}

func (r *gormRepo) SaveRevokedToken(ctx context.Context, jti string, exp int64) error {
	return r.db.WithContext(ctx).Create(&model.RevokedToken{JTI: jti, ExpiresAt: exp}).Error
}

func (r *gormRepo) ResetPassword(ctx context.Context, token, passwordHash string) error {
	ut, err := r.GetToken(ctx, token)
	if err != nil {
		return err
	}
	return r.db.WithContext(ctx).Model(&model.User{}).Where("email = ?", ut.Email).Update("password", passwordHash).Error
}
