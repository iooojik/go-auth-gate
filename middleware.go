package go_auth_gate

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/iooojik/go-auth-gate/internal/config"
	"github.com/iooojik/go-auth-gate/internal/middleware"
	"github.com/iooojik/go-auth-gate/internal/repository/session"
	"github.com/iooojik/go-auth-gate/internal/service/authservice"
	"github.com/iooojik/go-auth-gate/pkg/apple"
	"github.com/iooojik/go-auth-gate/pkg/google"
	"github.com/iooojik/go-auth-gate/pkg/jwt"
	"github.com/jmoiron/sqlx"
)

func NewMiddleware(ctx context.Context, cfgPath string) *middleware.Auth {
	cfg := config.Load(cfgPath)

	db, err := sqlx.ConnectContext(ctx, "mysql", cfg.SQL.SqlDsn)
	if err != nil {
		panic(err)
	}

	sessionsRepo := session.New(db)

	appleSecretsFile, err := os.Open(cfg.AppleSignIn.KeyPath)
	if err != nil {
		panic(fmt.Errorf("open apple_sign_in_key_file: %w", err))
	}

	defer func() { _ = appleSecretsFile.Close() }()

	appleSecretsContent, err := io.ReadAll(appleSecretsFile)
	if err != nil {
		panic(fmt.Errorf("read apple_sign_in_key_file: %w", err))
	}

	srv := authservice.New(
		apple.New(
			cfg.AppleSignIn,
			apple.GenerateClientSecret(appleSecretsContent),
			http.DefaultClient,
		),
		google.New(cfg.GoogleSignIn, http.DefaultClient),
		sessionsRepo,
	)

	authMiddleware := middleware.NewAuth(
		srv,
		jwt.TokenHeader,
		jwt.ValidateToken(cfg.JWT.SecretKey),
		jwt.GenerateToken(cfg.JWT.SecretKey, cfg.JWT.Domain),
	)

	return authMiddleware
}
