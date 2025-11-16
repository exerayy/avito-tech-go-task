package storage

import (
	"context"
	"database/sql"
)

type DB interface {
	Exec(ctx context.Context, query string, args ...any) (sql.Result, error)
	Query(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	Begin(ctx context.Context) (*sql.Tx, error)
}
