package session

import (
	"time"
)

type User struct {
	ID              int       `db:"id"`
	AppleID         string    `db:"apple_id"`
	SessionDuration int       `db:"session_duration"`
	CreatedAt       time.Time `db:"created_at"`
	AuthType        string    `db:"auth_type"`
}

type UserToken struct {
	ID        int       `db:"id"`
	UserID    int       `db:"user_id"`
	Token     string    `db:"token"`
	CreatedAt time.Time `db:"created_at"`
}
