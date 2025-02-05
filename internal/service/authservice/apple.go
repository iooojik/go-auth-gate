package authservice

import (
	"context"
	"fmt"

	"github.com/iooojik/go-auth-gate/internal/model"
	"github.com/iooojik/go-auth-gate/pkg/apple"
)

func (s *Service) AppleLogin(ctx context.Context, callbackInfo model.Generate) error {
	code, err := s.appleSignIn.ReceiveToken(ctx, apple.Generate{
		Code:   callbackInfo.Code,
		UserID: callbackInfo.UserID,
	})
	if err != nil {
		return fmt.Errorf("receive token %w", err)
	}

	err = s.sessionsRepository.Login(ctx, model.LoginInfo{
		UserID:         callbackInfo.UserID,
		AppleTokenInfo: code,
	})
	if err != nil {
		return fmt.Errorf("login: %w", err)
	}

	return nil
}

func (s *Service) AppleRefresh(ctx context.Context, refreshToken model.Refresh) error {
	code, err := s.appleSignIn.RefreshToken(ctx, apple.Refresh{
		RefreshToken: refreshToken.RefreshToken,
	})
	if err != nil {
		return fmt.Errorf("receive token %w", err)
	}

	err = s.sessionsRepository.Login(ctx, model.LoginInfo{
		UserID:         refreshToken.UserID,
		AppleTokenInfo: code,
	})
	if err != nil {
		return fmt.Errorf("login: %w", err)
	}

	return nil
}
