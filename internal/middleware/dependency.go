package middleware

import (
	"context"
)

type (
	SessionService interface {
		CheckUser(ctx context.Context, userID string) error
	}
)
