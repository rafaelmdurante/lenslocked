package models

import (
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"fmt"
	"github.com/rafaelmdurante/lenslocked/rand"
)

const (
	// MinBytesPerToken is the minimum number of bytes to be used for each session token
	MinBytesPerToken = 32
)

type Session struct {
	ID     int
	UserID int
	// Token is only set when creating a new session. When looking up a session
	// this will be left empty, as we only store the hash of a session token
	// in our database, and we cannot reverse it into a raw token.
	Token     string
	TokenHash string
}

type SessionService struct {
	DB *sql.DB
	// BytesPerToken is used to determine how many bytes to use when generating
	// each session token. If this value is not set or is less than the
	// MinBytesPerToken const it will be ignored and MinBytesPerToken will be
	// used instead.
	BytesPerToken int
}

// Create will create a new session for the user provided. The session token
// will be returned as the Token field on the Session type, but only the hashed
// session token is stored in the database.
func (ss *SessionService) Create(userID int) (*Session, error) {
	bytesPerToken := ss.BytesPerToken
	if bytesPerToken < MinBytesPerToken {
		bytesPerToken = MinBytesPerToken
	}

	token, err := rand.String(bytesPerToken)
	if err != nil {
		return nil, fmt.Errorf("error creating session token: %w", err)
	}

	session := Session{
		UserID:    userID,
		Token:     token,
		TokenHash: ss.hash(token),
	}

	// this 'on conflict ... do' is a POSTGRES specifically command
	row := ss.DB.QueryRow(`
        INSERT INTO sessions (user_id, token_hash)
            VALUES ($1, $2)
        ON CONFLICT (user_id) DO
        UPDATE
            SET token_hash = $2
        RETURNING id;`, session.UserID, session.TokenHash)
	err = row.Scan(&session.ID)

	// If the error was not sql.ErrNoRows, we need to check to see if it was any
	// other error. If it was sql.ErrNoRows it will be overwritten inside the if
	// block, and we still need to check for any errors.
	if err != nil {
		return nil, fmt.Errorf("error creating session token: %w", err)
	}

	return &session, nil
}

// User will return the user associated with that token. The tradeoff here is
// that the SessionService needs to know about the 'users' table and how to
// construct a User struct. It is a bit of intermingling responsibility though.
func (ss *SessionService) User(token string) (*User, error) {
	// 1. hash the session token
	tokenHash := ss.hash(token)
	var user User

	// 2. query for the session with that hash
	row := ss.DB.QueryRow(`
		SELECT
			u.id,
			u.email,
			u.password_hash
		FROM sessions s
			JOIN users u ON s.user_id = u.id
		WHERE s.token_hash = $1`, tokenHash)

	// 3. assign values to struct
	err := row.Scan(&user.ID, &user.Email, &user.PasswordHash)
	if err != nil {
		return nil, fmt.Errorf("user: error getting user %w", err)
	}

	// 4. return the user
	return &user, nil
}

func (ss *SessionService) hash(token string) string {
	tokenHash := sha256.Sum256([]byte(token))
	// base64 encode the data into a string
	return base64.URLEncoding.EncodeToString(tokenHash[:]) // [:] shorthand to [0:len(tokenHash)]
}

func (ss *SessionService) Delete(token string) error {
	tokenHash := ss.hash(token)

	// using Exec instead of QueryRow because we don't care for a return value
	_, err := ss.DB.Exec(`
    DELETE FROM sessions
    WHERE token_hash = $1;`, tokenHash)

	if err != nil {
		return fmt.Errorf("error deleting a session token: %w", err)
	}

	return nil
}
