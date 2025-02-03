package session

import (
	"context"
	"fmt"

	"github.com/iooojik/go-auth-gate/internal/model"
)

func (r *Repository) Login(ctx context.Context, loginInfo model.LoginInfo) error {
	// SQL-запрос
	query := `
		INSERT INTO users (apple_id, session_duration, created_at, auth_type)
		SELECT ?, ?, NOW(), 'apple_sign_in'
		WHERE NOT EXISTS (
			SELECT 1 FROM users WHERE apple_id = ?
		);
	
		INSERT INTO user_tokens (user_id, token, created_at)
		SELECT id, ?, NOW()
		FROM users WHERE apple_id = ?
		ON DUPLICATE KEY UPDATE token = VALUES(token), created_at = NOW();`

	// 1. Первый запрос: создаем пользователя и первый токен
	_, err := db.Exec(query, appleID, sessionDuration, appleID, token1, appleID)
	if err != nil {
		return fmt.Errorf("exec query: %w", err)
	}

	return nil
}
