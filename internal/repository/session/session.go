package session

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"iter"

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
		`INSERT INTO users (user_id, created_at) VALUES (?, NOW());`,
		loginInfo.UserID,
	)
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("insert user: %w", err)
	}

	if loginInfo.TokenType() == model.AppleID {
		err = r.insertAppleIDToken(ctx, tx, loginInfo)
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

func (r *Repository) insertAppleIDToken(ctx context.Context, tx *sql.Tx, loginInfo model.LoginInfo) error {
	ti := loginInfo.AppleTokenInfo

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
		loginInfo.UserID, ti.AccessToken, ti.TokenType, ti.ExpiresIn, ti.RefreshToken, ti.IDToken)
	if err != nil {
		return fmt.Errorf("insert tokens: %w", err)
	}

	return nil
}

func (r *Repository) CheckSession(_ context.Context, userID string) (bool, error) {
	var user User

	err := r.client.Get(&user, "SELECT id, user_id, created_at FROM users WHERE user_id = ?", userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, service.ErrUserDoesNotExists
		}

		return false, fmt.Errorf("get user: %w", err)
	}

	var token UserToken

	err = r.client.Get(&token, "SELECT id WHERE user_id = ?", userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}

		return false, fmt.Errorf("get token: %w", err)
	}

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
	rows, err := r.client.QueryContext(ctx, `
SELECT user_id,refresh_token FROM apple_tokens WHERE created_at >= NOW() - INTERVAL '30 minutes'`)
	if err != nil {
		return nil, fmt.Errorf("find user by id: %w", err)
	}

	defer func() { _ = rows.Close() }()

	return func(yield func(model.Refresh, error) bool) {
		for rows.Next() {
			var token model.Refresh

			err = rows.Scan(&token.UserID, &token.RefreshToken)
			if err != nil {
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
