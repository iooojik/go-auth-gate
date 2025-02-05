package authservice

import (
	"context"

	"github.com/iooojik/go-auth-gate/internal/model"
	"github.com/iooojik/go-auth-gate/pkg/apple"
)

type (
	AppleSignIn interface {
		ReceiveToken(ctx context.Context, gen apple.Generate) (*apple.AuthCode, error)

		RefreshToken(ctx context.Context, refresh apple.Refresh) (*apple.AuthCode, error)
	}

	GoogleSignIn interface {
		CheckToken(ctx context.Context, token string) (bool, error)
	}

	SessionRepository interface {
		Login(ctx context.Context, loginInfo model.LoginInfo) error

		CheckSession(_ context.Context, userID string) (bool, error)
	}
)
