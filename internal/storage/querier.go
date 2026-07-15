package storage

import (
	"context"
	"database/sql"
)

// Row is the minimal result contract needed by repositories for single-row queries.
type Row interface {
	Scan(dest ...any) error
}

// Querier is the SQL surface available to repositories and transaction callbacks.
type Querier interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) Row
}
