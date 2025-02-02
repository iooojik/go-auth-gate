package apple

import (
	"errors"
)

var (
	ErrBadRequest   = errors.New("bad request")
	ErrTokenRevoked = errors.New("token revoked")
)
