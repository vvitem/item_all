package sqlite

import (
	"context"
	"database/sql"
	"sync"
	"sync/atomic"
	"time"

	"github.com/vvitem/item_all/internal/storage"
	"github.com/vvitem/item_all/internal/storage/sqlite/migrations"
)

// Store owns the local SQLite connection pool and write serialization boundary.
type Store struct {
	db      *sql.DB
	path    string
	config  Config
	writeMu sync.Mutex
	closed  atomic.Bool
}

// Open initializes, verifies, migrates, and returns a ready Store.
func Open(ctx context.Context, config Config) (*Store, error) {
	if err := config.validate(); err != nil {
		return nil, err
	}
	path, err := prepareDatabasePath(config.Path)
	if err != nil {
		return nil, err
	}
	config.Path = path
	dsn, err := buildDSN(config)
	if err != nil {
		return nil, err
	}
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, classifySQLiteError("open-sqlite", err, storage.CodeOpenFailed, false)
	}
	fail := func(err error) (*Store, error) {
		_ = db.Close()
		return nil, err
	}
	db.SetMaxOpenConns(config.MaxOpenConns)
	db.SetMaxIdleConns(config.MaxIdleConns)
	db.SetConnMaxLifetime(0)
	if err := db.PingContext(ctx); err != nil {
		return fail(classifySQLiteError("ping-sqlite", err, storage.CodeOpenFailed, false))
	}
	if err := verifyPool(ctx, db, config); err != nil {
		return fail(err)
	}
	if err := runMigrations(ctx, db, migrations.FS, time.Now); err != nil {
		return fail(err)
	}
	if err := verifyPool(ctx, db, config); err != nil {
		return fail(err)
	}
	if err := secureDatabaseFile(path); err != nil {
		return fail(err)
	}
	return &Store{db: db, path: path, config: config}, nil
}

// Path returns the resolved local database file path.
func (s *Store) Path() string {
	if s == nil {
		return ""
	}
	return s.path
}

// ExecContext executes a serialized application write operation.
func (s *Store) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	if err := s.closedError("execute-sql"); err != nil {
		return nil, err
	}
	s.writeMu.Lock()
	defer s.writeMu.Unlock()
	if err := s.closedError("execute-sql"); err != nil {
		return nil, err
	}
	result, err := s.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, classifySQLiteError("execute-sql", err, storage.CodeTransactionFailed, false)
	}
	return result, nil
}

// QueryContext executes a read operation without taking the process-local write lock.
func (s *Store) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	if err := s.closedError("query-sql"); err != nil {
		return nil, err
	}
	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, classifySQLiteError("query-sql", err, storage.CodeTransactionFailed, false)
	}
	return rows, nil
}

// QueryRowContext returns a row adapter that classifies deferred Scan errors.
func (s *Store) QueryRowContext(ctx context.Context, query string, args ...any) storage.Row {
	if err := s.closedError("query-row-sql"); err != nil {
		return errorRow{err: err}
	}
	return sqlRow{row: s.db.QueryRowContext(ctx, query, args...), operation: "query-row-sql"}
}

// Close is idempotent and prevents subsequent operations from touching database/sql.
func (s *Store) Close() error {
	if s == nil || s.closed.Swap(true) {
		return nil
	}
	if err := s.db.Close(); err != nil {
		return classifySQLiteError("close-sqlite", err, storage.CodeOpenFailed, false)
	}
	return nil
}

func (s *Store) closedError(operation string) error {
	if s == nil || s.closed.Load() {
		return storage.Wrap(storage.CodeClosed, operation, sql.ErrConnDone, false)
	}
	return nil
}

type sqlRow struct {
	row       *sql.Row
	operation string
}

func (r sqlRow) Scan(dest ...any) error {
	if err := r.row.Scan(dest...); err != nil {
		return classifySQLiteError(r.operation, err, storage.CodeTransactionFailed, false)
	}
	return nil
}

type errorRow struct {
	err error
}

func (r errorRow) Scan(...any) error {
	return r.err
}
