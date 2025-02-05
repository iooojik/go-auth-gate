package service

import (
	"errors"
)

var (
	ErrUserDoesNotExists    = errors.New("user does not exists")
	ErrSessionDoesNotExists = errors.New("session does not exists")
)
