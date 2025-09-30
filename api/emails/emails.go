package emails

import (
	"context"
	"fmt"
	"os"

	gomail "gopkg.in/mail.v2"
)

type EmailsAPI interface {
	SendEmail(ctx context.Context, inp SendEmailInput) error
}

type SendEmailInput struct {
	To          string
	Subject     string
	ContentType string
	Body        string
}

type Emails struct{}

func New() *Emails {
	return &Emails{}
}

func (e *Emails) SendEmail(ctx context.Context, inp SendEmailInput) error {
	m := gomail.NewMessage()
	m.SetHeaders(map[string][]string{
		"From":    {os.Getenv("MAIL_USERNAME")},
		"To":      {inp.To},
		"Subject": {inp.Subject},
	})
	m.SetBody(inp.ContentType, inp.Body)

	d := gomail.NewDialer(os.Getenv("MAIL_HOST"), 587, os.Getenv("MAIL_USERNAME"), os.Getenv("MAIL_APP_PASSWORD"))
	if err := d.DialAndSend(m); err != nil {
		fmt.Println("ERROR", err)
		return err
	}

	return nil
}
