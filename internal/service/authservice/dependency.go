package authservice

import (
	"context"
	"iter"

	"github.com/iooojik/go-auth-gate/apple"
	"github.com/iooojik/go-auth-gate/internal/model"
)

type (
	AppleSignIn interface {
		ReceiveToken(ctx context.Context, gen apple.Generate) (*apple.AuthCode, error)

		RefreshToken(ctx context.Context, refresh apple.Refresh) (*apple.AuthCode, error)
	}

	GoogleSignIn interface {
		CheckToken(ctx context.Context, token string) (string, error)
	}

	SessionRepository interface {
		Login(ctx context.Context, loginInfo model.LoginInfo) error

		CheckSession(_ context.Context, userID string) (bool, error)

		FetchAll(ctx context.Context, authType model.TokenType) (iter.Seq2[model.Refresh, error], error)
	}
)
