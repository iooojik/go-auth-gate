package apple_test

import (
	"errors"
	"testing"

	"github.com/iooojik/go-auth-gate/pkg/apple"
	"github.com/stretchr/testify/require"
)

func TestGenerateClientSecret(t *testing.T) {
	t.Parallel()

	type args struct {
		privateKey []byte
		cfg        apple.TokenConfig
	}

	tests := []struct {
		name      string
		args      args
		wantToken string
		wantErr   error
	}{
		// TODO: Add test cases.
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotFunc := apple.GenerateClientSecret(tt.args.privateKey)

			token, err := gotFunc(tt.args.cfg)

			require.Equal(t, tt.wantToken, token)
			require.True(t, errors.Is(err, tt.wantErr))
		})
	}
}
