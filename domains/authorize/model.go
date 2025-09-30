package authorize

import (
	"github.com/msskobelina/fit-profi/pkg/mysql"
)

// swagger:model user
type User struct {
	ID       int    `json:"id,omitempty" gorm:"primaryKey"`
	FullName string `json:"fullName,omitempty"`
	Email    string `json:"email,omitempty" gorm:"not null;unique;index"`
	Password string `json:"password,omitempty"`

	mysql.Model
}

// swagger:model userToken
type UserToken struct {
	ID    int    `json:"id,omitempty" gorm:"primaryKey"`
	Email string `json:"email,omitempty" gorm:"index"`
	Token string `json:"token,omitempty"`

	mysql.Model
}

// swagger:model revokedToken
type RevokedToken struct {
	ID        int    `json:"id,omitempty" gorm:"primaryKey"`
	JTI       string `json:"jti,omitempty" gorm:"uniqueIndex;size:64"`
	ExpiresAt int64  `json:"expiresAt,omitempty"`

	mysql.Model
}

// swagger:model registerUserRequest
type RegisterUserRequest struct {
	FullName string `json:"fullName"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// swagger:model authResponse
type AuthResponse struct {
	Token    string `json:"token"`
	UserID   int    `json:"userId"`
	FullName string `json:"fullName"`
	Email    string `json:"email"`
	Error    *Error `json:"error,omitempty"`
}

// swagger:parameters usersRegister
type RegisterParams struct {
	// in: body
	// required: true
	Body RegisterUserRequest
}

// swagger:model loginUserRequest
type LoginUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// swagger:parameters usersLogin
type LoginParams struct {
	// in: body
	// required: true
	Body LoginUserRequest
}

// swagger:model sendEmailRequest
type SendEmailRequest struct {
	Email string `json:"email"`
}

// swagger:parameters usersSendEmail
type SendEmailParams struct {
	// in: body
	// required: true
	Body SendEmailRequest
}

// swagger:model resetPasswordRequest
type ResetPasswordRequest struct {
	Token    string `json:"token"`
	Password string `json:"password"`
}

// swagger:parameters usersResetPassword
type ResetPasswordParams struct {
	// in: body
	// required: true
	Body ResetPasswordRequest
}

// swagger:model serviceError
type Error struct {
	Code    string `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
	Status  int    `json:"status,omitempty"`
}

// swagger:model errorResponse
type ErrorResponse struct {
	// example: message
	Error string `json:"error"`
}
