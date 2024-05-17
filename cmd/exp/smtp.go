package main

import (
	"fmt"

	"github.com/rafaelmdurante/lenslocked/models"
)

// all values from the mail service - in this case, mailtrap
const (
	host     = "sandbox.smtp.mailtrap.io"
	port     = 587
	username = "c193f039338455"
	password = "7118eaa1f8d5eb"
)

// main function re-declared in the package because it is experimental
// ignore LSP error/warnings
func main() {
	email := models.Email{
		From:      "test@lenslocked.com",
		To:        "raf@ael.com",
		Subject:   "this is a test email",
		Plaintext: "this is the body of the email",
		HTML:      `<h1>Hello there!</h1><p>This is the HTML email</p>`,
	}

    es := models.NewEmailService(models.SMTPConfig{
        Host: host,
        Port: port,
        Username: username,
        Password: password,
    })

    err := es.Send(email)
    if err != nil {
        panic(err)
    }
    fmt.Println("Email sent")
}
