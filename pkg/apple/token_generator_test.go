package apple_test

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/iooojik/go-auth-gate/pkg/apple"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func generateTestPrivateKey(t *testing.T) []byte {
	privKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	require.NoError(t, err)

	derKey, err := x509.MarshalPKCS8PrivateKey(privKey)
	require.NoError(t, err)

	return pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: derKey})
}

func TestGenerateClientSecret(t *testing.T) {
	t.Parallel()

	type args struct {
		privateKey []byte
		cfg        apple.TokenConfig
	}

	tests := []struct {
		name     string
		args     args
		validate func(t *testing.T, a args, token string, err error)
	}{
		{
			name: "test#1",
			args: args{
				privateKey: nil,
				cfg: apple.TokenConfig{
					TeamID:   "test-team",
					ClientID: "test-Client",
					KeyID:    "test-key",
					Audience: "https://example.com",
					ExpSec:   3600,
				},
			},
			validate: func(t *testing.T, a args, token string, err error) {
				require.NoError(t, err)

				parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (any, error) {
					claims := token.Claims.(jwt.MapClaims)

					require.Equal(t, "https://example.com", claims["aud"])
					require.Equal(t, "test-team", claims["iss"])
					require.Equal(t, "test-Client", claims["sub"])

					iat := time.Unix(int64(claims["iat"].(float64)), 0)
					exp := time.Unix(int64(claims["exp"].(float64)), 0)

					require.True(t, iat.Before(exp))
					require.True(t, exp.Sub(iat).Seconds() == 3600)

					return apple.ParseECPrivateKey(a.privateKey).Public(), nil
				})

				require.NoError(t, err)
				assert.True(t, parsedToken.Valid)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			tt.args.privateKey = generateTestPrivateKey(t)

			gotFunc := apple.GenerateClientSecret(tt.args.privateKey)

			token, err := gotFunc(tt.args.cfg)

			tt.validate(t, tt.args, token, err)
		})
	}
}
