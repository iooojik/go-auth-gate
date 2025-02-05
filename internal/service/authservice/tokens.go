package authservice

import (
	"context"
	"iter"

	"github.com/iooojik/go-auth-gate/internal/model"
)

func (s *Service) FetchAll(
	ctx context.Context, authType model.TokenType,
) (iter.Seq2[model.Refresh, error], error) {
	return s.sessionsRepository.FetchAll(ctx, authType)
}
