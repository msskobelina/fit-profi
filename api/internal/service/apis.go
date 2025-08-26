package service

import "context"

type APIs struct {
	Emails EmailsAPI
}

type EmailsAPI interface {
	SendEmail(ctx context.Context, inp SendEmailInput) error
}

type SendEmailInput struct {
	To          string
	Subject     string
	ContentType string
	Body        string
}
