package models

import "github.com/go-mail/mail/v2"

const (
	// DefaultSender is the default email address to send emails from
	DefaultSender = "support@lenslocked.com"
)

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

