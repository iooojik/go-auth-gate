package apple

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type ClientSecretGenerator func(cfg TokenConfig) (string, error)

type TokenConfig struct {
	ClientID string `yaml:"clientID"`
	TeamID   string `yaml:"teamID"`
	KeyID    string `yaml:"keyID"`
	Audience string `yaml:"audience"`
	ExpSec   int    `yaml:"exp"`
}

func GenerateClientSecret(privateKey []byte) ClientSecretGenerator {
	ecdsaKey := ParseECPrivateKey(privateKey)

	return func(cfg TokenConfig) (string, error) {
		headers := jwt.MapClaims{
			"iss": cfg.TeamID,
			"iat": time.Now().Unix(),
			"exp": time.Now().Add(time.Duration(cfg.ExpSec) * time.Second).Unix(),
			"aud": cfg.Audience,
			"sub": cfg.ClientID,
		}

		token := jwt.NewWithClaims(jwt.SigningMethodES256, headers)
		token.Header["kid"] = cfg.KeyID

		signedToken, err := token.SignedString(ecdsaKey)
		if err != nil {
			return "", fmt.Errorf("signing token: %w", err)
		}

		return signedToken, nil
	}
}

func ParseECPrivateKey(pemKey []byte) *ecdsa.PrivateKey {
	block, _ := pem.Decode(pemKey)
	if block == nil || block.Type != "PRIVATE KEY" {
		panic(ErrDecodePEMBlock)
	}

	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		panic(fmt.Errorf("failed to parse private key: %w", err))
	}

	ecdsaKey, ok := key.(*ecdsa.PrivateKey)
	if !ok {
		panic(ErrKeyIsNotECDSA)
	}

	return ecdsaKey
}
