package apple

import (
	"errors"
)

var (
	ErrBadRequest     = errors.New("bad request")
	ErrTokenRevoked   = errors.New("token revoked")
	ErrKeyIsNotECDSA  = errors.New("key is not of type ECDSA")
	ErrDecodePEMBlock = errors.New("decode PEM block containing private key")
	ErrEmptyLink      = errors.New("empty link")
)
