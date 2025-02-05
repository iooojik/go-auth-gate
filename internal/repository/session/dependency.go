package session

import (
	"database/sql"
)

type (
	Conn interface {
		Begin() (*sql.Tx, error)

		Get(dest any, query string, args ...any) error

		// QueryContext(ctx context.Context, query string, insertArgs ...any) (*sql.Rows, error)

		// ExecContext(ctx context.Context, query string, insertArgs ...any) (sql.Result, error)
	}
)
