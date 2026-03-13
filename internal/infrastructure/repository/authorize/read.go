package authorize

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/msskobelina/fit-profi/internal/domain/model"
	utilsErrors "github.com/msskobelina/fit-profi/pkg/errors"
)

func (r *gormRepo) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	var u model.User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&u).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &utilsErrors.Error{Message: "User not registered"}
		}
		return nil, err
	}
	return &u, nil
}

func (r *gormRepo) GetToken(ctx context.Context, token string) (*model.UserToken, error) {
	var ut model.UserToken
	err := r.db.WithContext(ctx).Where("token = ?", token).First(&ut).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &utilsErrors.Error{Message: "User with provided token not found"}
		}
		return nil, err
	}
	return &ut, nil
}

func (r *gormRepo) IsTokenRevoked(ctx context.Context, jti string) (bool, error) {
	var n int64
	if err := r.db.WithContext(ctx).Model(&model.RevokedToken{}).Where("jti = ?", jti).Count(&n).Error; err != nil {
		return false, err
	}
	return n > 0, nil
}
