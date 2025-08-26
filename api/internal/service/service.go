package service

import (
	"context"
	"encoding/json"
	"fmt"
)

///
/// ERRORS
///

type Error struct {
	Code    string `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
	Status  int    `json:"status,omitempty"`
}

func (e *Error) Error() string {
	b, err := json.Marshal(e)
	if err != nil {
		return fmt.Sprintf("failed to marshal error: %v", err)
	}
	return string(b)
}

///
/// SERVICES
///

type Services struct {
	Users Users
}

type Users interface {
	RegisterUser(ctx context.Context, inp *RegisterUserInput) (RegisterUserOutput, error)
	LoginUser(ctx context.Context, inp *LoginUserInput) (LoginUserOutput, error)
	Logout(ctx context.Context, token string) error
	VerifyAccessToken(ctx context.Context, token string) (bool, int, string)
	SendEmail(ctx context.Context, inp *SendUserEmailInput) error
	ResetPassword(ctx context.Context, inp *ResetPasswordInput) error
}

type RegisterUserInput struct {
	FullName string `json:"fullName"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
type RegisterUserOutput struct {
	Token    string
	UserID   int
	FullName string
	Email    string
}

type LoginUserInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type LoginUserOutput struct {
	Token    string
	UserID   int
	FullName string
	Email    string
}

type SendUserEmailInput struct {
	Email string `json:"email"`
	Token string `json:"token"`
}
