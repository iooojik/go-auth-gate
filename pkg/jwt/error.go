package jwt

import (
	"errors"
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrEmptyUserID  = errors.New("empty user ID")
)
