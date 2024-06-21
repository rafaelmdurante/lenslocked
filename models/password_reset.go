package models

import (
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/rafaelmdurante/lenslocked/rand"
)

const (
	// DefaultResetDuration is the default time that a PasswordReset is valid for
	DefaultResetDuration = 1 * time.Hour
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

func (prs *PasswordResetService) Create(email string) (*PasswordReset, error) {
	// verify we have a valid email address for a user
	email = strings.ToLower(email)

	var userID int
	row := prs.DB.QueryRow(`
		SELECT id FROM users WHERE email = $1;`, email)

	err := row.Scan(&userID)
	if err != nil {
		// TODO: consider returning a specific error when the user does not exist
		return nil, fmt.Errorf("create: %w", err)
	}

	// build the password reset
	bytesPerToken := prs.BytesPerToken
	if bytesPerToken == 0 {
		bytesPerToken = MinBytesPerToken
	}

	token, err := rand.String(bytesPerToken)
	if err != nil {
		return nil, fmt.Errorf("create: %w", err)
	}

	duration := prs.Duration
	if duration == 0 {
		duration = DefaultResetDuration
	}

	pwReset := PasswordReset{
		UserID: userID,
		Token: token,
		TokenHash: prs.hash(token),
		ExpiresAt: time.Now().Add(duration),
	}

	row = prs.DB.QueryRow(`
		INSERT INTO password_resets (user_id, token_hash, expires_at)
		VALUES ($1, $2, $3) ON CONFLICT (user_id) DO
		UPDATE
		SET token_hash = $2, expires_at = $3
		RETURNING id;`, pwReset.UserID, pwReset.TokenHash, pwReset.ExpiresAt)

	err = row.Scan(&pwReset.ID)
	if err != nil {
		return nil, fmt.Errorf("create: %w", err)
	}

	return &pwReset, nil
}

// We are going to consume a token and return the user associated with it,
// or return an error if the token wasn't valid for any reason
func (prs *PasswordResetService) Consume(token string) (*User, error) {
	return nil, fmt.Errorf("TODO: Implement PasswordResetService.Consume")
}


func (prs *PasswordResetService) hash(token string) string {
	tokenHash := sha256.Sum256([]byte(token))

	return base64.URLEncoding.EncodeToString(tokenHash[:])
}
