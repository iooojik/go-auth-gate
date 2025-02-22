package session

import (
	"context"
	"database/sql"
)

type (
	Conn interface {
		Begin() (*sql.Tx, error)

		Get(dest any, query string, args ...any) error

		QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	}
)
