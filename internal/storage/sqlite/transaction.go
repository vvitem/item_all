package sqlite

import (
	"context"
	"database/sql"
	"errors"

	"github.com/vvitem/item_all/internal/storage"
)

type transactionContextKey struct{}

func transactionContextError(ctx context.Context, operation string) error {
	if ctx != nil && ctx.Value(transactionContextKey{}) != nil {
		return storage.Wrap(storage.CodeTransactionFailed, operation, errors.New("use the transaction Querier inside a transaction callback"), false)
	}
	return nil
}

type rawTransaction interface {
	ExecContext(context.Context, string, ...any) (sql.Result, error)
	QueryContext(context.Context, string, ...any) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...any) *sql.Row
	Commit() error
	Rollback() error
}

// WithTx serializes application writes and owns commit or rollback.
func (s *Store) WithTx(ctx context.Context, fn func(context.Context, storage.Querier) error) error {
	if err := s.closedError("begin-transaction"); err != nil {
		return err
	}
	if fn == nil {
		return storage.Wrap(storage.CodeTransactionFailed, "validate-transaction-callback", errors.New("transaction callback is required"), false)
	}
	if ctx.Value(transactionContextKey{}) != nil {
		return storage.Wrap(storage.CodeTransactionFailed, "reject-nested-transaction", errors.New("nested transactions are unsupported"), false)
	}

	s.writeMu.Lock()
	defer s.writeMu.Unlock()
	if err := s.closedError("begin-transaction"); err != nil {
		return err
	}
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return classifySQLiteError("begin-transaction", err, storage.CodeTransactionFailed, false)
	}
	txCtx := context.WithValue(ctx, transactionContextKey{}, struct{}{})
	return runTransaction(txCtx, tx, fn)
}

func runTransaction(ctx context.Context, tx rawTransaction, fn func(context.Context, storage.Querier) error) (err error) {
	if tx == nil || fn == nil {
		return storage.Wrap(storage.CodeTransactionFailed, "validate-transaction", errors.New("transaction and callback are required"), false)
	}
	defer func() {
		if recovered := recover(); recovered != nil {
			_ = tx.Rollback()
			panic(recovered)
		}
	}()

	querier := transactionQuerier{tx: tx}
	if callbackErr := fn(ctx, querier); callbackErr != nil {
		return rollbackTransaction(tx, callbackErr)
	}
	if contextErr := ctx.Err(); contextErr != nil {
		classified := storage.Wrap(storage.CodeTransactionFailed, "transaction-context", contextErr, false)
		return rollbackTransaction(tx, classified)
	}
	if commitErr := tx.Commit(); commitErr != nil {
		return storage.Wrap(storage.CodeCommitFailed, "commit-transaction", commitErr, false)
	}
	return nil
}

func rollbackTransaction(tx rawTransaction, cause error) error {
	rollbackErr := tx.Rollback()
	if rollbackErr == nil || errors.Is(rollbackErr, sql.ErrTxDone) {
		return cause
	}
	classified := storage.Wrap(storage.CodeTransactionFailed, "rollback-transaction", rollbackErr, false)
	return errors.Join(cause, classified)
}

type transactionQuerier struct {
	tx rawTransaction
}

func (q transactionQuerier) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	result, err := q.tx.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, classifySQLiteError("execute-transaction-sql", err, storage.CodeTransactionFailed, false)
	}
	return result, nil
}

func (q transactionQuerier) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	rows, err := q.tx.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, classifySQLiteError("query-transaction-sql", err, storage.CodeTransactionFailed, false)
	}
	return rows, nil
}

func (q transactionQuerier) QueryRowContext(ctx context.Context, query string, args ...any) storage.Row {
	return sqlRow{row: q.tx.QueryRowContext(ctx, query, args...), operation: "query-row-transaction-sql"}
}
