package apple_test

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/iooojik/go-auth-gate/apple"
	"github.com/iooojik/go-auth-gate/apple/mocks"
	"github.com/stretchr/testify/require"
)

func TestClient_RefreshToken(t *testing.T) {
	t.Parallel()

	type fields struct {
		cfg            apple.Config
		tokenGenerator apple.ClientSecretGenerator
	}

	type args struct {
		httpClient *mocks.HTTPClient
		refresh    apple.Refresh
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		setup   func(ctx context.Context, a *args)
		want    *apple.AuthCode
		wantErr bool
	}{
		{
			name: "test#1",
			fields: fields{
				tokenGenerator: func(_ apple.TokenConfig) (string, error) {
					return "clientToken1", nil
				},
				//nolint:exhaustruct
				cfg: apple.Config{
					URL: "https://appleid.apple.com",
				},
			},
			//nolint:exhaustruct
			args: args{
				refresh: apple.Refresh{
					RefreshToken: "refreshToken",
				},
			},
			setup: func(ctx context.Context, a *args) {
				//nolint:revive
				u, err := url.Parse("https://appleid.apple.com/auth/token?client_id=&client_secret=clientToken1&grant_type=refresh_token&refresh_token=refreshToken")
				require.NoError(t, err)

				req, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), nil)
				require.NoError(t, err)

				//nolint:exhaustruct
				a.httpClient.
					EXPECT().
					Do(req).
					Return(&http.Response{
						StatusCode: http.StatusOK,
						Body: io.NopCloser(strings.NewReader(
							"{\n    \"access_token\": \"access_token_body\"," +
								"\n    \"token_type\": \"Bearer\",\n    \"expires_in\": 881," +
								"\n    \"refresh_token\": \"refresh_token_body\"," +
								"\n    \"id_token\": \"id_token_body\"\n}")),
					}, nil)
			},
			want: &apple.AuthCode{
				AccessToken:  "access_token_body",
				TokenType:    "Bearer",
				ExpiresIn:    881,
				RefreshToken: "refresh_token_body",
				IDToken:      "id_token_body",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			tt.args.httpClient = mocks.NewHTTPClient(t)

			ctx := context.Background()

			tt.setup(ctx, &tt.args)

			r := apple.New(tt.fields.cfg, tt.fields.tokenGenerator, tt.args.httpClient)

			got, err := r.RefreshToken(ctx, tt.args.refresh)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, tt.want, got)
		})
	}
}

func TestClient_ReceiveToken(t *testing.T) {
	t.Parallel()

	type fields struct {
		cfg            apple.Config
		tokenGenerator apple.ClientSecretGenerator
	}

	type args struct {
		httpClient *mocks.HTTPClient
		gen        apple.Generate
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		setup   func(ctx context.Context, a *args)
		want    *apple.AuthCode
		wantErr bool
	}{
		{
			name: "test#1",
			fields: fields{
				tokenGenerator: func(_ apple.TokenConfig) (string, error) {
					return "clientToken1", nil
				},
				//nolint:exhaustruct
				cfg: apple.Config{
					URL: "https://appleid.apple.com",
				},
			},
			//nolint:exhaustruct
			args: args{
				gen: apple.Generate{
					Code: "123",
				},
			},
			setup: func(ctx context.Context, a *args) {
				//nolint:revive
				u, err := url.Parse("https://appleid.apple.com/auth/token?client_id=&client_secret=clientToken1&code=123&grant_type=authorization_code")
				require.NoError(t, err)

				req, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), nil)
				require.NoError(t, err)

				//nolint:exhaustruct
				a.httpClient.
					EXPECT().
					Do(req).
					Return(&http.Response{
						StatusCode: http.StatusOK,
						Body: io.NopCloser(strings.NewReader(
							"{\n    \"access_token\": \"access_token_body\"," +
								"\n    \"token_type\": \"Bearer\",\n    \"expires_in\": 881," +
								"\n    \"refresh_token\": \"refresh_token_body\"," +
								"\n    \"id_token\": \"id_token_body\"\n}")),
					}, nil)
			},
			want: &apple.AuthCode{
				AccessToken:  "access_token_body",
				TokenType:    "Bearer",
				ExpiresIn:    881,
				RefreshToken: "refresh_token_body",
				IDToken:      "id_token_body",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			tt.args.httpClient = mocks.NewHTTPClient(t)

			ctx := context.Background()

			tt.setup(ctx, &tt.args)

			r := apple.New(tt.fields.cfg, tt.fields.tokenGenerator, tt.args.httpClient)

			got, err := r.ReceiveToken(ctx, tt.args.gen)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, tt.want, got)
		})
	}
}
