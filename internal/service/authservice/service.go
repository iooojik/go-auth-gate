package authservice

type Service struct {
	appleSignIn        AppleSignIn
	googleSignIn       GoogleSignIn
	sessionsRepository SessionRepository
}

func New(
	appleSignIn AppleSignIn,
	googleSignIn GoogleSignIn,
	sessionsRepository SessionRepository,
) *Service {
	s := &Service{
		appleSignIn:        appleSignIn,
		googleSignIn:       googleSignIn,
		sessionsRepository: sessionsRepository,
	}

	return s
}
