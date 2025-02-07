package apple

import (
	"context"
	"fmt"

	"github.com/iooojik/go-auth-gate/internal/model"
)

type Refresh struct {
	sessionService SessionService
}

func New(sessionService SessionService) *Refresh {
	r := &Refresh{
		sessionService: sessionService,
	}

	return r
}

func (r *Refresh) Run(ctx context.Context) error {
	iterSeq, err := r.sessionService.FetchAll(ctx, model.AppleID)
	if err != nil {
		return fmt.Errorf("fetch apple tokens: %w", err)
	}

	for token, err := range iterSeq {
		if err != nil {
			return fmt.Errorf("fetch apple tokens: %w", err)
		}

		err = r.sessionService.AppleRefresh(ctx, token)
		if err != nil {
			return fmt.Errorf("refresh apple tokens: %w, userID: %s", err, token.UserID)
		}
	}

	return nil
}
