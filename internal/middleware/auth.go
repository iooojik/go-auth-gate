package middleware

import (
	"net/http"

	"github.com/iooojik/go-auth-gate/pkg/jwt"
)

func Login(
	srv SessionService,
	tokenHeader string,
	validator jwt.TokenValidator,
	generator jwt.TokenGenerator,
) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// 1. receive auth data from request.
			// 2. determine token type.
			// 3. login.

			next.ServeHTTP(w, r)
		})
	}
}

func Auth(
	srv SessionService,
	tokenHeader string,
	validator jwt.TokenValidator,
	generator jwt.TokenGenerator,
) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// 1. receive token.
			token := r.Header.Get(tokenHeader)
			// 2. validate token.
			claims, err := validator(token)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)

				return
			}
			// 3. check user.
			err = srv.CheckUser(r.Context(), claims.TokenUser.ID)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)

				return
			}

			// 4. generate new token.
			token, err = generator(claims.TokenUser.ID)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)

				return
			}

			w.Header().Add(tokenHeader, token)

			next.ServeHTTP(w, r)
		})
	}
}
