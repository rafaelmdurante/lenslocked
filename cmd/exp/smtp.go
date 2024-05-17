package main

import (
	"os"

	"github.com/go-mail/mail/v2"
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

    dialer := mail.NewDialer(host, port, username, password)
    
    // two ways of doing so, the first is shorter, for single emails
    err := dialer.DialAndSend(msg)
    if err != nil {
        // TODO: handle error correctly
        panic(err)
    }

    // second one is handy when you're looping and sending multiple emails
    sender, err := dialer.Dial()
    if err != nil {
        // TODO: handle error correctly
        panic(err)
    }
    defer sender.Close()

    err = mail.Send(sender, msg)
    if err != nil {
        // TODO: handle the error correctly
        panic(err)
    }
}

