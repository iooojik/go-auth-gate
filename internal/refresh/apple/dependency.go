package apple

import (
	"context"
	"iter"

	"github.com/iooojik/go-auth-gate/internal/model"
)

type (
	SessionService interface {
		FetchAll(
			ctx context.Context, authType model.TokenType,
		) (iter.Seq2[model.Refresh, error], error)

		AppleRefresh(ctx context.Context, refreshToken model.Refresh) error
	}
)
