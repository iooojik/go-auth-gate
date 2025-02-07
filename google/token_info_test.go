package google_test

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/iooojik/go-auth-gate/pkg/google"
	"github.com/iooojik/go-auth-gate/pkg/google/mocks"
	"github.com/stretchr/testify/require"
)

func TestClient_TokenInfo(t *testing.T) {
	t.Parallel()

	type fields struct {
		cfg google.Config
	}

	type args struct {
		httpClient *mocks.HTTPClient
		token      string
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		setup   func(ctx context.Context, a *args)
		want    *google.TokenInfo
		wantErr bool
	}{
		{
			name: "test#1",
			fields: fields{
				//nolint:exhaustruct
				cfg: google.Config{
					URL: "https://oauth2.googleapis.com",
				},
			},
			//nolint:exhaustruct
			args: args{
				token: "123456",
			},
			setup: func(ctx context.Context, a *args) {
				u, err := url.Parse("https://oauth2.googleapis.com/tokeninfo?id_token=123456")
				require.NoError(t, err)

				req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
				require.NoError(t, err)

				a.httpClient.
					EXPECT().
					Do(req).
					Return(
						//nolint:exhaustruct
						&http.Response{
							StatusCode: http.StatusOK,
							Body:       io.NopCloser(strings.NewReader("{}")),
						}, nil)
			},
			//nolint:exhaustruct
			want:    &google.TokenInfo{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			tt.args.httpClient = mocks.NewHTTPClient(t)

			ctx := context.Background()

			tt.setup(ctx, &tt.args)

			r := google.New(tt.fields.cfg, tt.args.httpClient)

			got, err := r.TokenInfo(ctx, tt.args.token)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, tt.want, got)
		})
	}
}
