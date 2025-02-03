package service

// func (s *Service) AppleCallback(ctx context.Context, callbackInfo model.AppleCallback) error {
// 	code, err := s.appleSignIn.ReceiveToken(ctx, apple.Generate{Code: callbackInfo.Code})
// 	if err != nil {
// 		return fmt.Errorf("receive token %w", err)
// 	}
//
// 	err = s.userRepo.Login(ctx, model.LoginInfo{}) // todo add user info + token info.
// 	if err != nil {
// 		return fmt.Errorf("login: %w", err)
// 	}
//
// 	return nil
// }
//
// func (s *Service) RefreshToken(ctx context.Context) error {
// 	return nil
// }
