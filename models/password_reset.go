package models

import (
	"database/sql"
	"fmt"
	"time"
)

const (
	// DefaultResetDuration is the default time that a PasswordReset is valid for
	DefaultResetDurante = 1 * time.Hour
)

type PasswordReset struct {
	ID     int
	UserID int
	// Token is only set when a PasswordReset is being created
	// this is because we only store the hash in our database
	// and we cannot recreate it when looking up a password reset in our db
	Token     string
	TokenHash string
	ExpiresAt time.Time
}

type PasswordResetService struct {
	DB *sql.DB
	// BytesPerToken is used to determine how many bytes to use when generating
	// each password reset token. If this value is not set or is less than the
	// MiinBytesPerToken const it will be ignore and MinBytesPerToken will be
	// used.
	BytesPerToken int
	// Duration is the amount of time that a PasswordReset is valid for.
	// Defaults to DefaultResetDuration
	Duration time.Duration
}

func (prs *PasswordResetService) Create(userID int) (*PasswordReset, error) {
	return nil, fmt.Errorf("TODO: implement PasswordResetService.Create")
}

// We are going to consume a token and return the user associated with it,
// or return an error if the token wasn't valid for any reason
func (prs *PasswordResetService) Consume(token string) (*User, error) {
	return nil, fmt.Errorf("TODO: Implement PasswordResetService.Consume")
}
