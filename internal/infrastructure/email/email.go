package email

import (
	"context"
	"fmt"
	"os"

	gomail "gopkg.in/mail.v2"
)

type SendInput struct {
	To          string
	Subject     string
	ContentType string
	Body        string
}

type Sender interface {
	Send(ctx context.Context, inp SendInput) error
}

type smtpSender struct{}

func NewSender() Sender {
	return &smtpSender{}
}

func (s *smtpSender) Send(_ context.Context, inp SendInput) error {
	m := gomail.NewMessage()
	m.SetHeaders(map[string][]string{
		"From":    {os.Getenv("MAIL_USERNAME")},
		"To":      {inp.To},
		"Subject": {inp.Subject},
	})
	m.SetBody(inp.ContentType, inp.Body)

	d := gomail.NewDialer(os.Getenv("MAIL_HOST"), 587, os.Getenv("MAIL_USERNAME"), os.Getenv("MAIL_APP_PASSWORD"))
	if err := d.DialAndSend(m); err != nil {
		fmt.Println("mail send error:", err)
		return err
	}

	return nil
}
