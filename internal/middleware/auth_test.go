package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/iooojik/go-auth-gate/internal/middleware"
	"github.com/iooojik/go-auth-gate/internal/middleware/mocks"
	"github.com/iooojik/go-auth-gate/internal/model"
	"github.com/iooojik/go-auth-gate/internal/service"
	"github.com/iooojik/go-auth-gate/pkg/jwt"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestAuth_Auth(t *testing.T) {
	t.Parallel()

	type fields struct {
		srv         *mocks.SessionService
		tokenHeader string
		validator   jwt.TokenValidator
		generator   jwt.TokenGenerator
	}

	type args struct {
		next http.Handler
		req  *http.Request
	}

	tests := []struct {
		name           string
		fields         fields
		setup          func(f *fields, a *args)
		args           args
		wantToken      bool
		wantStatusCode int
	}{
		{
			name: "test#1",
			fields: fields{
				srv:         nil,
				tokenHeader: jwt.TokenHeader,
				validator:   jwt.ValidateToken("123"),
				generator:   jwt.GenerateToken("123", "example.com"),
			},
			setup: func(f *fields, a *args) {
				f.srv.
					EXPECT().
					CheckUser(mock.Anything, "88").
					Return(nil)

				testToken, err := f.generator("88")
				require.NoError(t, err)

				a.req.Header.Set(f.tokenHeader, testToken)
			},
			args: args{
				req: httptest.NewRequest(http.MethodGet, "/", nil),
				next: http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
					w.WriteHeader(http.StatusOK)
				}),
			},
			wantToken:      true,
			wantStatusCode: http.StatusOK,
		},
		{
			name: "test#2",
			fields: fields{
				srv:         nil,
				tokenHeader: jwt.TokenHeader,
				validator:   jwt.ValidateToken("123"),
				generator:   jwt.GenerateToken("123", "example.com"),
			},
			setup: func(f *fields, a *args) {
				f.srv.
					EXPECT().
					CheckUser(mock.Anything, "88").
					Return(service.ErrUserDoesNotExists)

				testToken, err := f.generator("88")
				require.NoError(t, err)

				a.req.Header.Set(f.tokenHeader, testToken)
			},
			args: args{
				req: httptest.NewRequest(http.MethodGet, "/", nil),
				next: http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
					w.WriteHeader(http.StatusOK)
				}),
			},
			wantToken:      false,
			wantStatusCode: http.StatusUnauthorized,
		},
	}

	//nolint:dupl
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			tt.fields.srv = mocks.NewSessionService(t)

			tt.setup(&tt.fields, &tt.args)

			a := middleware.NewAuth(
				tt.fields.srv,
				tt.fields.tokenHeader,
				tt.fields.validator,
				tt.fields.generator,
			)

			got := a.Auth(tt.args.next)

			httpTest := httptest.NewRecorder()

			got.ServeHTTP(httpTest, tt.args.req)

			if tt.wantToken {
				require.NotEmpty(t, httpTest.Header().Get(tt.fields.tokenHeader))
			} else {
				require.Empty(t, httpTest.Header().Get(tt.fields.tokenHeader))
			}

			require.Equal(t, tt.wantStatusCode, httpTest.Code)
		})
	}
}

func TestAuth_Login(t *testing.T) {
	t.Parallel()

	type fields struct {
		srv         *mocks.SessionService
		tokenHeader string
		validator   jwt.TokenValidator
		generator   jwt.TokenGenerator
	}

	type args struct {
		next http.Handler
		req  *http.Request
	}

	tests := []struct {
		name           string
		fields         fields
		setup          func(f *fields, a *args)
		args           args
		wantToken      bool
		wantStatusCode int
	}{
		{
			name: "test#1",
			fields: fields{
				srv:         nil,
				tokenHeader: jwt.TokenHeader,
				validator:   jwt.ValidateToken("123"),
				generator:   jwt.GenerateToken("123", "example.com"),
			},
			setup: func(_ *fields, _ *args) {

			},
			args: args{
				req: httptest.NewRequest(http.MethodGet, "https://example.com", nil),
				next: http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
					w.WriteHeader(http.StatusOK)
				}),
			},
			wantToken:      false,
			wantStatusCode: http.StatusForbidden,
		},
		{
			name: "test#2",
			fields: fields{
				srv:         nil,
				tokenHeader: jwt.TokenHeader,
				validator:   jwt.ValidateToken("123"),
				generator:   jwt.GenerateToken("123", "example.com"),
			},
			setup: func(f *fields, _ *args) {
				f.srv.EXPECT().
					GoogleLogin(mock.Anything, "google_token").
					Return("user_id1", nil)
			},
			args: args{
				req: httptest.NewRequest(http.MethodGet, "https://example.com?token=google_token", nil),
				next: http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
					w.WriteHeader(http.StatusOK)
				}),
			},
			wantToken:      true,
			wantStatusCode: http.StatusOK,
		},
		{
			name: "test#3",
			fields: fields{
				srv:         nil,
				tokenHeader: jwt.TokenHeader,
				validator:   jwt.ValidateToken("123"),
				generator:   jwt.GenerateToken("123", "example.com"),
			},
			setup: func(f *fields, _ *args) {
				f.srv.EXPECT().
					AppleLogin(mock.Anything, model.Generate{
						Code:   "apple_code",
						UserID: "99",
					}).
					Return(nil)
			},
			args: args{
				req: httptest.NewRequest(
					http.MethodGet,
					"https://example.com?code=apple_code&client_id=99",
					nil,
				),
				next: http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
					w.WriteHeader(http.StatusOK)
				}),
			},
			wantToken:      true,
			wantStatusCode: http.StatusOK,
		},
	}

	//nolint:dupl
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			tt.fields.srv = mocks.NewSessionService(t)

			tt.setup(&tt.fields, &tt.args)

			a := middleware.NewAuth(
				tt.fields.srv,
				tt.fields.tokenHeader,
				tt.fields.validator,
				tt.fields.generator,
			)

			got := a.Login(tt.args.next)

			httpTest := httptest.NewRecorder()

			got.ServeHTTP(httpTest, tt.args.req)

			if tt.wantToken {
				require.NotEmpty(t, httpTest.Header().Get(tt.fields.tokenHeader))
			} else {
				require.Empty(t, httpTest.Header().Get(tt.fields.tokenHeader))
			}

			require.Equal(t, tt.wantStatusCode, httpTest.Code)
		})
	}
}
