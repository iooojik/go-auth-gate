package authmiddleware

import (
	"github.com/iooojik/go-auth-gate/internal/model"
)

type AuthData struct {
	ClientID string
	Code     string
	Token    string
}

func (a *AuthData) AuthType() model.TokenType {
	if a.ClientID != "" && a.Code != "" {
		return model.AppleID
	}

	if a.Token != "" {
		return model.GoogleSignInAuth
	}

	return model.Unknown
}
