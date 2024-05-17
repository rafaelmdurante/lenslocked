package models

import (
	"fmt"

	"github.com/go-mail/mail/v2"
)

const (
	// DefaultSender is the default email address to send emails from
	DefaultSender = "support@lenslocked.com"
)

type Email struct {
	From      string
	To        string
	Subject   string
	Plaintext string
	HTML      string
}

type EmailService struct {
	// DefaultSender is used as the default sender when one isn't provided for
	// an email. This is also used in function where the email is a
	// predetermined, like the forgotten password email.
	DefaultSender string

	dialer *mail.Dialer
}

type SMTPConfig struct {
	Host     string
	Port     int
	Username string
	Password string
}

func NewEmailService(config SMTPConfig) *EmailService {
	return &EmailService{
		dialer: mail.NewDialer(
			config.Host, config.Port, config.Username, config.Password),
	}
}

func (es *EmailService) Send(email Email) error {
	msg := mail.NewMessage()

	msg.SetHeader("To", email.To)
    es.setFrom(msg, email)
	msg.SetHeader("From", email.From)
	msg.SetHeader("Subject", email.Subject)

	switch {
	case email.Plaintext != "" && email.HTML != "":
		msg.SetBody("text/plain", email.Plaintext)
		msg.AddAlternative("text/html", email.HTML)
	case email.Plaintext != "":
		msg.SetBody("text/plain", email.Plaintext)
	case email.HTML != "":
		msg.SetBody("text/html", email.HTML)
	}

	err := es.dialer.DialAndSend(msg)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

func (es *EmailService) setFrom(msg *mail.Message, email Email) {
	var from string

	switch {
    case email.From != "":
        from = email.From
    case es.DefaultSender != "":
        from = es.DefaultSender
    default:
        from = DefaultSender
	}

    msg.SetHeader("From", from)
}
