package jwt

import (
	"fmt"
	"time"

	json "github.com/json-iterator/go"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const (
	TokenHeader = "Token"
)

type (
	TokenGenerator func(id string) (string, error)

	TokenValidator func(token string) (*TokenClaims, error)
)

type TokenClaims struct {
	TokenUser `json:"user"`
	jwt.RegisteredClaims
}

func ValidateToken(secret string) TokenValidator {
	return func(headerToken string) (*TokenClaims, error) {
		if headerToken == "" {
			return nil, nil
		}

		token, err := jwt.Parse(headerToken, func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, ErrInvalidToken
			}

			return []byte(secret), nil
		})
		if err != nil {
			return new(TokenClaims), err
		}

		mappedData, _ := token.Claims.(jwt.MapClaims)

		data, err := json.Marshal(mappedData)
		if err != nil {
			return nil, fmt.Errorf("marshal claims: %w", err)
		}

		ctx := new(TokenClaims)

		err = json.Unmarshal(data, ctx)
		if err != nil {
			return nil, fmt.Errorf("unmarshal claims: %w", err)
		}

		return ctx, nil
	}
}

func GenerateToken(key, domain string) TokenGenerator {
	return func(id string) (string, error) {
		if id == "" {
			return "", ErrEmptyUserID
		}

		claims := TokenClaims{
			TokenUser: TokenUser{
				ID: id,
			},
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer:    domain,
				Subject:   id,
				Audience:  []string{"apps"},
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(14 * 24 * time.Hour)),
				NotBefore: jwt.NewNumericDate(time.Now().Add(-10 * time.Minute)),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
				ID:        uuid.New().String(),
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		ss, err := token.SignedString([]byte(key))
		if err != nil {
			return "", fmt.Errorf("sign token: %w", err)
		}

		return ss, nil
	}
}
