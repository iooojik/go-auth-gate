package google

import (
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_checkResponse(t *testing.T) {
	t.Parallel()

	type args struct {
		resp *http.Response
	}

	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name: "test#1",
			args: args{
				resp: &http.Response{
					StatusCode: http.StatusOK,
				},
			},
			wantErr: nil,
		},
		{
			name: "test#2",
			args: args{
				resp: &http.Response{
					StatusCode: http.StatusBadRequest,
					Body:       io.NopCloser(strings.NewReader("abcdef")),
				},
			},
			wantErr: ErrBadRequest,
		},
		{
			name: "test#3",
			args: args{
				resp: &http.Response{
					StatusCode: http.StatusBadRequest,
					Body:       io.NopCloser(strings.NewReader("{\n    \"error\": \"invalid_token\",\n    \"error_description\": \"Invalid Value\"\n}")),
				},
			},
			wantErr: ErrTokenRevokedOrInvalid,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := checkResponse(tt.args.resp)

			if tt.wantErr == nil {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}

func Test_prepareLink(t *testing.T) {
	t.Parallel()

	type args struct {
		link   string
		params url.Values
	}

	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "test#1",
			args: args{
				link: "https://google.com",
				params: url.Values{
					"key": []string{"value"},
				},
			},
			want:    "https://google.com/tokeninfo?key=value",
			wantErr: false,
		},
		{
			name: "test#2",
			args: args{
				link: "https://google.com/",
				params: url.Values{
					"key": []string{"value"},
				},
			},
			want:    "https://google.com/tokeninfo?key=value",
			wantErr: false,
		},
		{
			name: "test#3",
			args: args{
				link:   "https://google.com/",
				params: url.Values{},
			},
			want:    "https://google.com/tokeninfo",
			wantErr: false,
		},
		{
			name: "test#4",
			args: args{
				link:   "",
				params: url.Values{},
			},
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := prepareLink(tt.args.link, tt.args.params)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			gotLink := ""

			if got != nil {
				gotLink = got.String()
			}

			require.Equal(t, tt.want, gotLink)
		})
	}
}
