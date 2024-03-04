// Package rand is a tailored wrap for the crypto/rand native Go package.
// It is aimed to help development of our application by providing functions
// specially tailored fow our own needs.
package rand

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

const SessionTokenBytes = 32

// SessionToken is a helper function that generates a random session token
func SessionToken() (string, error) {
	return String(SessionTokenBytes)
}

// String returns a random string using crypto/rand.
// n is the number of bytes being used to generate the random string.
func String(n int) (string, error) {
	b, err := Bytes(n)
	if err != nil {
		return "", fmt.Errorf("string: %w", err)
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// Bytes function returns a slice of bytes of random data to be used for
// randomness
func Bytes(n int) ([]byte, error) {
	b := make([]byte, n)

	// reads a slice of bytes and return the amount of bytes successfully read
	nRead, err := rand.Read(b)
	if err != nil {
		return nil, fmt.Errorf("bytes: %w", err)
	}

	// the amount of bytes must be the same as the size of the slice
	if nRead < n {
		return nil, fmt.Errorf("bytes: didn't read enough random bytes")
	}

	return b, nil
}
