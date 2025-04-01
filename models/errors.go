package models

import "errors"

// a common pattern is to add the package as a prefix to the error for context
var (
	ErrEmailToken = errors.New("models: email address is already in use")
	ErrNotFound   = errors.New("models: resource could not be found")
)
