package authservice

import (
	"context"
	"fmt"
	"iter"

	"github.com/iooojik/go-auth-gate/internal/model"
)

func (s *Service) FetchAll(
	ctx context.Context, authType model.TokenType,
) (iter.Seq2[model.Refresh, error], error) {
	seq, err := s.sessionsRepository.FetchAll(ctx, authType)
	if err != nil {
		return nil, fmt.Errorf("fetch all tokens: %w", err)
	}

	return seq, nil
}
