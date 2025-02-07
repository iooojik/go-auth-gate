package authmiddleware

import (
	"github.com/iooojik/go-auth-gate/pkg/jwt"
)

type Auth struct {
	srv         SessionService
	tokenHeader string
	validator   jwt.TokenValidator
	generator   jwt.TokenGenerator
}

func NewAuth(
	srv SessionService,
	tokenHeader string,
	validator jwt.TokenValidator,
	generator jwt.TokenGenerator,
) *Auth {
	a := &Auth{
		srv:         srv,
		tokenHeader: tokenHeader,
		validator:   validator,
		generator:   generator,
	}

	return a
}
