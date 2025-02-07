package authservice

import (
	"context"
	"fmt"

	"github.com/iooojik/go-auth-gate/internal/service"
)

func (s *Service) CheckUser(ctx context.Context, userID string) error {
	exists, err := s.sessionsRepository.CheckSession(ctx, userID)
	if err != nil {
		return fmt.Errorf("check user: %w", err)
	}

	if !exists {
		return service.ErrSessionDoesNotExists
	}

	return nil
}
