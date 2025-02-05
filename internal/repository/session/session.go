package session

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

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

	_, err := tx.ExecContext(ctx, `INSERT INTO user_tokens 
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
