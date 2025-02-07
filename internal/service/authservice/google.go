package authservice

import (
	"context"
	"fmt"

	"github.com/iooojik/go-auth-gate/internal/model"
	"github.com/iooojik/go-auth-gate/internal/service"
)

func (s *Service) GoogleLogin(ctx context.Context, token string) (string, error) {
	userID, err := s.googleSignIn.CheckToken(ctx, token)
	if err != nil {
		return "", fmt.Errorf("check goolge token %w", err)
	}

	err = s.sessionsRepository.Login(ctx, model.LoginInfo{
		UserID:         userID,
		AppleTokenInfo: nil,
	})
	if err != nil {
		return "", fmt.Errorf("login: %w", err)
	}

	if userID == "" {
		return "", fmt.Errorf("%w: %s", service.ErrInvalidToken, token)
	}

	return userID, nil
}
