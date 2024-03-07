package models

import (
	"database/sql"
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

	// TODO: hash the session token

	session := Session{
		UserID: userID,
		Token:  token,
		// TODO: set the TokenHash
	}

	// TODO: store the session in our DB

	return &session, nil
}

// User will return the user associated with that token. The tradeoff here is
// that the SessionService needs to know about the 'users' table and how to
// construct a User struct. It is a bit of intermingling responsibility though.
func (ss *SessionService) User(token string) (*User, error) {
	// TODO: implement ss.User
	return nil, nil
}
