package service

import (
	"context"
	"fmt"

	"github.com/iooojik/go-auth-gate/internal/model"
	"github.com/iooojik/go-auth-gate/pkg/apple"
)

func (s *Service) AppleCallback(ctx context.Context, callbackInfo apple.Generate) error {
	code, err := s.appleSignIn.ReceiveToken(ctx, callbackInfo)
	if err != nil {
		return fmt.Errorf("receive token %w", err)
	}

	err = s.sessionsRepository.Login(ctx, model.LoginInfo{
		UserID:    callbackInfo.UserID,
		Token:     code.RefreshToken,
		TokenType: "apple_sign_in",
	}) // todo save token from apple.
	if err != nil {
		return fmt.Errorf("login: %w", err)
	}

	return nil
}
