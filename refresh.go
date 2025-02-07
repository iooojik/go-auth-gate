package authgate

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/iooojik/go-auth-gate/internal/config"
	applerefresh "github.com/iooojik/go-auth-gate/internal/refresh/apple"
	"github.com/iooojik/go-auth-gate/internal/repository/session"
	"github.com/iooojik/go-auth-gate/internal/service/authservice"
	"github.com/iooojik/go-auth-gate/pkg/apple"
	"github.com/iooojik/go-auth-gate/pkg/google"
	"github.com/jmoiron/sqlx"
)

func RunRefresh(ctx context.Context, cfg config.Config) error {
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

	r := applerefresh.New(srv)

	err = r.Run(ctx)
	if err != nil {
		return fmt.Errorf("run refresh: %w", err)
	}

	return nil
}
