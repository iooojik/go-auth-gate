package session

import (
	"time"

	"github.com/iooojik/go-auth-gate/internal/model"
)

type User struct {
	ID       int             `db:"id"`
	UserID   string          `db:"user_id"`
	AuthType model.TokenType `db:"auth_type"`
}

type UserToken struct {
	ID           int       `db:"id"`
	UserID       string    `db:"user_id"`
	AccessToken  string    `db:"access_token"`
	TokenType    string    `db:"token_type"`
	ExpiresIn    int       `db:"expires_in"`
	RefreshToken string    `db:"refresh_token"`
	IDToken      string    `db:"id_token"`
	CreatedAt    time.Time `db:"created_at"`
}
