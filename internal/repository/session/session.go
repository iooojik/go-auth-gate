package session

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"iter"

	"github.com/iooojik/go-auth-gate/apple"
	"github.com/iooojik/go-auth-gate/internal/model"
	"github.com/iooojik/go-auth-gate/internal/service"
)

func (r *Repository) Login(ctx context.Context, loginInfo model.LoginInfo) error {
	tx, err := r.client.Begin()
	if err != nil {
		return fmt.Errorf("start transaction: %w", err)
	}

	_, err = tx.ExecContext(
		ctx,
		`INSERT INTO users (user_id, created_at, auth_type) VALUES (?, NOW(), ?);`,
		loginInfo.UserID,
		loginInfo.TokenType(),
	)
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("insert user: %w", err)
	}

	if loginInfo.TokenType() == model.AppleID {
		err = InsertAppleIDToken(ctx, tx, *loginInfo.AppleTokenInfo, loginInfo.UserID)
		if err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("insert apple id token: %w", err)
		}
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("commit transacion: %w", err)
	}

	return nil
}

func InsertAppleIDToken(
	ctx context.Context,
	tx *sql.Tx,
	ti apple.AuthCode,
	userID string,
) error {
	_, err := tx.ExecContext(ctx, `INSERT INTO apple_tokens 
    (user_id, access_token, token_type, expires_in, refresh_token, id_token, created_at)
VALUES (?, ?, ?, ?, ?, ?, NOW())
ON DUPLICATE KEY UPDATE 
    access_token = VALUES(access_token), 
    token_type = VALUES(token_type), 
    expires_in = VALUES(expires_in), 
    refresh_token = VALUES(refresh_token), 
    id_token = VALUES(id_token), 
	created_at = NOW();`,
		userID, ti.AccessToken, ti.TokenType, ti.ExpiresIn, ti.RefreshToken, ti.IDToken)
	if err != nil {
		return fmt.Errorf("insert tokens: %w", err)
	}

	return nil
}

func (r *Repository) CheckSession(_ context.Context, userID string) (bool, error) {
	var user User

	err := r.client.Get(&user, "SELECT user_id, auth_type  FROM users WHERE user_id = ?", userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, service.ErrUserDoesNotExists
		}

		return false, fmt.Errorf("get user: %w", err)
	}

	if user.AuthType == model.GoogleSignInAuth {
		return user.UserID != "", nil
	}

	var token UserToken

	err = r.client.Get(&token, "SELECT id FROM apple_tokens WHERE user_id = ?", userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}

		return false, fmt.Errorf("get token: %w", err)
	}
	// 1740140498
	return user.UserID != "" && token.IDToken != "", nil
}

func (r *Repository) FetchAll(
	ctx context.Context, authType model.TokenType,
) (iter.Seq2[model.Refresh, error], error) {
	if authType == model.AppleID {
		return r.fetchAppleTokens(ctx)
	}

	return nil, service.ErrUnknownAuthMethod
}

func (r *Repository) fetchAppleTokens(ctx context.Context) (iter.Seq2[model.Refresh, error], error) {
	rows, err := r.client.QueryContext(ctx, `SELECT user_id, refresh_token 
FROM apple_tokens WHERE created_at >= NOW() - INTERVAL 30 MINUTE;`)
	if err != nil {
		return nil, fmt.Errorf("find tokens by id: %w", err)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("find tokens by id: %w", rows.Err())
	}

	return func(yield func(model.Refresh, error) bool) {
		defer func() { _ = rows.Close() }()

		for rows.Next() {
			var token model.Refresh

			err = rows.Scan(&token.UserID, &token.RefreshToken)
			if err != nil {
				//nolint:exhaustruct
				yield(model.Refresh{}, err)

				return
			}

			more := yield(token, nil)
			if !more {
				return
			}
		}
	}, nil
}
