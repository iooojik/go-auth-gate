package service

type Service struct {
	appleSignIn AppleSignIn
	userRepo    UserRepository
}

func New() *Service {
	s := &Service{}

	return s
}
