package apple

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
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
	ecdsaKey, err := parseECPrivateKey(privateKey)
	if err != nil {
		panic(fmt.Errorf("parsing private key: %w", err))
	}

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

func parseECPrivateKey(pemKey []byte) (*ecdsa.PrivateKey, error) {
	block, _ := pem.Decode(pemKey)
	if block == nil || block.Type != "PRIVATE KEY" {
		return nil, errors.New("failed to decode PEM block containing private key")
	}

	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	ecdsaKey, ok := key.(*ecdsa.PrivateKey)
	if !ok {
		return nil, errors.New("key is not of type ECDSA")
	}

	return ecdsaKey, nil
}
