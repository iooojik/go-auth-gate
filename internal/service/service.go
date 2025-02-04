package service

type Service struct {
	appleSignIn        AppleSignIn
	sessionsRepository SessionRepository
}

func New() *Service {
	s := &Service{}

	return s
}
