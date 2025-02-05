package authservice

type Service struct {
	appleSignIn        AppleSignIn
	googleSignIn       GoogleSignIn
	sessionsRepository SessionRepository
}

func New() *Service {
	s := &Service{}

	return s
}
