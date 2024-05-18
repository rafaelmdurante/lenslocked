package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/rafaelmdurante/lenslocked/models"
)

// main function re-declared in the package because it is experimental
// ignore LSP error/warnings
func main() {
    err := godotenv.Load()
    if err != nil {
        log.Fatal("error loading .env file")
    }

    host := os.Getenv("SMTP_HOST")
    port, err := strconv.Atoi(os.Getenv("SMTP_PORT"))
    if err != nil {
        panic(err)
    }
    username := os.Getenv("SMTP_USERNAME")
    password := os.Getenv("SMTP_PASSWORD")

    es := models.NewEmailService(models.SMTPConfig{
        Host: host,
        Port: port,
        Username: username,
        Password: password,
    })

    err = es.ForgotPassword("pipersyd@proton.me", "https://lenslocked.com/reset-pw?token=abc123")
    if err != nil {
        panic(err)
    }
    fmt.Println("Email sent")
}
