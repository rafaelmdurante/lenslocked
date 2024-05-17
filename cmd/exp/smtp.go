package main

import (
	"os"

	"github.com/go-mail/mail/v2"
)

// main function re-declared in the package because it is experimental
// ignore LSP error/warnings
func main() {
    from := "test@lenslocked.com"
    to := "raf@ael.com"
    subject := "this is a test email"
    plaintext := "this is the body of the email"
    html := `<h1>Hello there!</h1><p>This is the HTML email</p>`

    msg := mail.NewMessage()
    msg.SetHeader("To", to)
    msg.SetHeader("From", from)
    msg.SetHeader("Subject", subject)
    msg.SetBody("text/plain", plaintext)
    msg.AddAlternative("text/html", html)
    msg.WriteTo(os.Stdout)
}

