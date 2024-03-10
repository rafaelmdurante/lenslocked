package models

import (
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"errors"
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

	row := ss.DB.QueryRow(`
		UPDATE sessions
		SET token_hash = $2
		WHERE user_id = $1
		RETURNING id;`, session.UserID, session.TokenHash)
	err = row.Scan(&session.ID)

	if errors.Is(err, sql.ErrNoRows) {
		// If no sessions exist, we will get an error. That means we need to
		// create a session object for that user.
		row = ss.DB.QueryRow(`
		INSERT INTO sessions (user_id, token_hash)
		VALUES ($1, $2)
		RETURNING id;`, session.UserID, session.TokenHash)
		// The error will be overwritten with either a new error or nil
		err = row.Scan(&session.ID)
	}

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
	// TODO: implement ss.User
	return nil, nil
}

func (ss *SessionService) hash(token string) string {
	tokenHash := sha256.Sum256([]byte(token))
	// base64 encode the data into a string
	return base64.URLEncoding.EncodeToString(tokenHash[:]) // [:] shorthand to [0:len(tokenHash)]
}
