package authgate

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql" // compile driver.
	"github.com/iooojik/go-auth-gate/authmiddleware"
	"github.com/iooojik/go-auth-gate/internal/repository/session"
	"github.com/iooojik/go-auth-gate/internal/service/authservice"
	"github.com/iooojik/go-auth-gate/pkg/apple"
	"github.com/iooojik/go-auth-gate/pkg/google"
	"github.com/iooojik/go-auth-gate/pkg/jwt"
	"github.com/jmoiron/sqlx"
)

func NewMiddleware(ctx context.Context, cfg Config) *authmiddleware.Auth {
	db, err := sqlx.ConnectContext(ctx, "mysql", cfg.SQL.SQLDsn)
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

	authMiddleware := authmiddleware.NewAuth(
		srv,
		jwt.TokenHeader,
		jwt.ValidateToken(cfg.JWT.SecretKey),
		jwt.GenerateToken(cfg.JWT.SecretKey, cfg.JWT.Domain),
	)

	return authMiddleware
}
