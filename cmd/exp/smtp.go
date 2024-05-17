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
    es := models.NewEmailService(models.SMTPConfig{
        Host: host,
        Port: port,
        Username: username,
        Password: password,
    })

    err := es.ForgotPassword("pipersyd@proton.me", "https://lenslocked.com/reset-pw?token=abc123")
    if err != nil {
        panic(err)
    }
    fmt.Println("Email sent")
}
