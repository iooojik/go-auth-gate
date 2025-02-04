package session

import (
	"context"
	"fmt"

	"github.com/iooojik/go-auth-gate/internal/model"
)

func (r *Repository) Login(ctx context.Context, loginInfo model.LoginInfo) error {
	tx, err := r.client.Begin()
	if err != nil {
		return fmt.Errorf("start transaction: %w", err)
	}

	_, err = tx.ExecContext(ctx, `
INSERT INTO users (user_id, session_duration, created_at, auth_type)
VALUES (?, ?, NOW(), ?)
ON DUPLICATE KEY UPDATE 
    session_duration = VALUES(session_duration);
`, loginInfo.UserID, r.cfg.SessionDuration, loginInfo.TokenType)
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("insert user: %w", err)
	}

	_, err = tx.ExecContext(ctx, `INSERT INTO user_tokens (user_id, token, created_at)
VALUES (?, ?, NOW())
ON DUPLICATE KEY UPDATE token = VALUES(token), created_at = NOW();`,
		loginInfo.UserID, loginInfo.Token)
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("insert tokens: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("commit transacion: %w", err)
	}

	return nil
}
