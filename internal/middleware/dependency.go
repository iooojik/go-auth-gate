package middleware

import (
	"context"

	"github.com/iooojik/go-auth-gate/internal/model"
)

type (
	SessionService interface {
		CheckUser(ctx context.Context, userID string) error

		AppleLogin(ctx context.Context, callbackInfo model.Generate) error

		GoogleLogin(ctx context.Context, token string) (string, error)
	}
)
