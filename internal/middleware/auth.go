package middleware

import (
	"net/http"

	"github.com/iooojik/go-auth-gate/internal/model"
)

func (a *Auth) Login(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()

		// 1. receive auth data from request.
		authData := AuthData{
			ClientID: query.Get("client_id"),
			Code:     query.Get("code"),
			Token:    query.Get("token"),
		}
		// 2. determine token type.
		authType := authData.AuthType()

		// 3. login.
		var err error

		switch authType {
		case model.AppleID:
			err = a.srv.AppleLogin(r.Context(), model.Generate{
				Code:   authData.Code,
				UserID: authData.ClientID,
			})

		case model.GoogleSignInAuth:
			authData.ClientID, err = a.srv.GoogleLogin(r.Context(), authData.Token)

		case model.Unknown:
			w.WriteHeader(http.StatusForbidden)

			return

		default:
			w.WriteHeader(http.StatusForbidden)

			return
		}

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)

			return
		}

		// 4. generate new token.
		token, err := a.generator(authData.ClientID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)

			return
		}

		w.Header().Add(a.tokenHeader, token)

		next.ServeHTTP(w, r)
	})
}

func (a *Auth) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 1. receive token.
		token := r.Header.Get(a.tokenHeader)
		// 2. validate token.
		claims, err := a.validator(token)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)

			return
		}
		// 3. check user.
		err = a.srv.CheckUser(r.Context(), claims.TokenUser.ID)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)

			return
		}

		// 4. generate new token.
		token, err = a.generator(claims.TokenUser.ID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)

			return
		}

		w.Header().Add(a.tokenHeader, token)

		next.ServeHTTP(w, r)
	})
}
