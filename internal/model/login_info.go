package model

import (
	"github.com/iooojik/go-auth-gate/pkg/apple"
)

type TokenType int

const (
	AppleID TokenType = iota
	GoogleSignInAuth
)

func (t TokenType) String() string {
	switch t {
	case AppleID:
		return "AppleID"

	case GoogleSignInAuth:
		return "GoogleSignInAuth"

	default:
		panic("unknown token type")
	}
}

type LoginInfo struct {
	UserID         string
	AppleTokenInfo *apple.AuthCode
}

func (i *LoginInfo) TokenType() TokenType {
	if i.AppleTokenInfo != nil {
		return AppleID
	}

	return GoogleSignInAuth
}
