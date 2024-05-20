package models

import (
	"bytes"
	"fmt"
	"html/template"
	"strings"

	"github.com/go-mail/mail/v2"
	"github.com/rafaelmdurante/lenslocked/templates"
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

func (es *EmailService) ForgotPassword(to, resetURL string) error {
	email := Email{
		Subject: "Reset your password",
		To:      to,
		// Plaintext: "To reset your password, please visit the following link: " + resetURL,
		// HTML: `<p>To reset your password, please visit the following link: <a href="` + resetURL + `">` + resetURL,
	}

	// TEMPLATES GENERAL
	data := struct {
	    ResetURL string
	} {
	    ResetURL: resetURL,
	}

	var emailHtml = "email.html"

	// HTML
    html := template.New(emailHtml) // it has to be the filename: https://stackoverflow.com/a/49043639

        html, err := html.ParseFS(templates.FS, emailHtml)
	if err != nil {
		return fmt.Errorf("error parsing html: %w", err)
	}

	var b bytes.Buffer
	err = html.Execute(&b, data)
	if err != nil {
		return fmt.Errorf("error executing html: %w", err)
	}
	email.HTML = b.String()

	// PLAINTEXT
    var emailPlain = "email.plain"
	tmplPlain, err := template.New(emailPlain).ParseFS(templates.FS, emailPlain)
	if err != nil {
		return fmt.Errorf("error parsing plain: %w", err)
	}

    // another approach rather than bytes buffer
	var plainBuilder strings.Builder
	err = tmplPlain.Execute(&plainBuilder, data)
	if err != nil {
		return fmt.Errorf("error executing plain: %w", err)
	}
	email.Plaintext = plainBuilder.String()

    // double check
    fmt.Println(email)

	// SEND EMAIL
	err = es.Send(email)
	if err != nil {
		return fmt.Errorf("forgot password email: %w", err)
	}

	return nil
}
