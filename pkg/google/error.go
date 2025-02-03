package google

import (
	"errors"
)

var (
	ErrBadRequest            = errors.New("bad request")
	ErrTokenRevokedOrInvalid = errors.New("token revoked or invalid")
	ErrEmptyLink             = errors.New("empty link")
)
