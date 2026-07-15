package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/vvitem/item_all/internal/storage"
)

func TestWithTxCommitsOnSuccess(t *testing.T) {
	store := openTestStore(t)
	ctx := context.Background()
	if _, err := store.ExecContext(ctx, `CREATE TABLE tx_values(value TEXT NOT NULL) STRICT`); err != nil {
		t.Fatal(err)
	}
	if err := store.WithTx(ctx, func(txCtx context.Context, q storage.Querier) error {
		_, err := q.ExecContext(txCtx, `INSERT INTO tx_values(value) VALUES(?)`, "committed")
		return err
	}); err != nil {
		t.Fatal(err)
	}
	var count int
	if err := store.QueryRowContext(ctx, `SELECT COUNT(*) FROM tx_values WHERE value='committed'`).Scan(&count); err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Fatalf("count = %d", count)
	}
}

func TestWithTxRollsBackCallbackError(t *testing.T) {
	store := openTestStore(t)
	ctx := context.Background()
	if _, err := store.ExecContext(ctx, `CREATE TABLE tx_values(value TEXT NOT NULL) STRICT`); err != nil {
		t.Fatal(err)
	}
	callbackErr := errors.New("callback failed")
	err := store.WithTx(ctx, func(txCtx context.Context, q storage.Querier) error {
		if _, err := q.ExecContext(txCtx, `INSERT INTO tx_values(value) VALUES(?)`, "rolled-back"); err != nil {
			return err
		}
		return callbackErr
	})
	if !errors.Is(err, callbackErr) {
		t.Fatalf("error = %v", err)
	}
	var count int
	if err := store.QueryRowContext(ctx, `SELECT COUNT(*) FROM tx_values`).Scan(&count); err != nil {
		t.Fatal(err)
	}
	if count != 0 {
		t.Fatalf("count = %d", count)
	}
}

func TestWithTxRollsBackAndRepanics(t *testing.T) {
	store := openTestStore(t)
	ctx := context.Background()
	if _, err := store.ExecContext(ctx, `CREATE TABLE tx_values(value TEXT NOT NULL) STRICT`); err != nil {
		t.Fatal(err)
	}
	const panicValue = "panic-value"
	func() {
		defer func() {
			if got := recover(); got != panicValue {
				t.Fatalf("panic = %#v", got)
			}
		}()
		_ = store.WithTx(ctx, func(txCtx context.Context, q storage.Querier) error {
			if _, err := q.ExecContext(txCtx, `INSERT INTO tx_values(value) VALUES(?)`, "rolled-back"); err != nil {
				return err
			}
			panic(panicValue)
		})
	}()
	var count int
	if err := store.QueryRowContext(ctx, `SELECT COUNT(*) FROM tx_values`).Scan(&count); err != nil {
		t.Fatal(err)
	}
	if count != 0 {
		t.Fatalf("count = %d", count)
	}
}

func TestWithTxHonorsContextCancellation(t *testing.T) {
	store := openTestStore(t)
	ctx, cancel := context.WithCancel(context.Background())
	err := store.WithTx(ctx, func(txCtx context.Context, _ storage.Querier) error {
		cancel()
		return nil
	})
	if !storage.IsCode(err, storage.CodeTransactionFailed) || !errors.Is(err, context.Canceled) {
		t.Fatalf("error = %v", err)
	}
}

func TestWithTxRejectsNilCallback(t *testing.T) {
	store := openTestStore(t)
	if err := store.WithTx(context.Background(), nil); !storage.IsCode(err, storage.CodeTransactionFailed) {
		t.Fatalf("error = %v", err)
	}
}

func TestWithTxRejectsNestedTransaction(t *testing.T) {
	store := openTestStore(t)
	err := store.WithTx(context.Background(), func(txCtx context.Context, _ storage.Querier) error {
		return store.WithTx(txCtx, func(context.Context, storage.Querier) error { return nil })
	})
	if !storage.IsCode(err, storage.CodeTransactionFailed) {
		t.Fatalf("error = %v", err)
	}
}

func TestWithTxRejectsClosedStore(t *testing.T) {
	store := openTestStore(t)
	if err := store.Close(); err != nil {
		t.Fatal(err)
	}
	err := store.WithTx(context.Background(), func(context.Context, storage.Querier) error { return nil })
	if !storage.IsCode(err, storage.CodeClosed) {
		t.Fatalf("error = %v", err)
	}
}

func TestRunTransactionMapsCommitFailure(t *testing.T) {
	commitErr := errors.New("commit failed")
	tx := &fakeRawTransaction{commitErr: commitErr}
	err := runTransaction(context.Background(), tx, func(context.Context, storage.Querier) error { return nil })
	if !storage.IsCode(err, storage.CodeCommitFailed) || !errors.Is(err, commitErr) {
		t.Fatalf("error = %v", err)
	}
}

func TestRunTransactionPreservesCallbackAndRollbackErrors(t *testing.T) {
	callbackErr := errors.New("callback failed")
	rollbackErr := errors.New("rollback failed")
	tx := &fakeRawTransaction{rollbackErr: rollbackErr}
	err := runTransaction(context.Background(), tx, func(context.Context, storage.Querier) error { return callbackErr })
	if !errors.Is(err, callbackErr) || !errors.Is(err, rollbackErr) {
		t.Fatalf("error = %v", err)
	}
}

type fakeRawTransaction struct {
	commitErr   error
	rollbackErr error
}

func (*fakeRawTransaction) ExecContext(context.Context, string, ...any) (sql.Result, error) {
	return nil, nil
}
func (*fakeRawTransaction) QueryContext(context.Context, string, ...any) (*sql.Rows, error) {
	return nil, nil
}
func (*fakeRawTransaction) QueryRowContext(context.Context, string, ...any) *sql.Row { return nil }
func (f *fakeRawTransaction) Commit() error                                          { return f.commitErr }
func (f *fakeRawTransaction) Rollback() error                                        { return f.rollbackErr }

func TestStoreRejectsTransactionContextOutsideQuerier(t *testing.T) {
	store := openTestStore(t)
	ctx := context.WithValue(context.Background(), transactionContextKey{}, struct{}{})
	if _, err := store.ExecContext(ctx, `SELECT 1`); !storage.IsCode(err, storage.CodeTransactionFailed) {
		t.Fatalf("Exec error = %v", err)
	}
	if _, err := store.QueryContext(ctx, `SELECT 1`); !storage.IsCode(err, storage.CodeTransactionFailed) {
		t.Fatalf("Query error = %v", err)
	}
	if err := store.QueryRowContext(ctx, `SELECT 1`).Scan(new(int)); !storage.IsCode(err, storage.CodeTransactionFailed) {
		t.Fatalf("QueryRow error = %v", err)
	}
}
