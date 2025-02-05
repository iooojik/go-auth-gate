package authservice

import (
	"context"
	"fmt"

	"github.com/iooojik/go-auth-gate/internal/model"
	"github.com/iooojik/go-auth-gate/pkg/apple"
)

func (s *Service) AppleLogin(ctx context.Context, callbackInfo apple.Generate) error {
	code, err := s.appleSignIn.ReceiveToken(ctx, callbackInfo)
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
