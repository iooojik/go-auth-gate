package service

import (
	"errors"
)

var (
	ErrUserDoesNotExists    = errors.New("user does not exists")
	ErrSessionDoesNotExists = errors.New("session does not exists")
	ErrUnknownAuthMethod    = errors.New("unknown auth method")
	ErrInvalidToken         = errors.New("invalid token")
)
