package sqlite

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"github.com/vvitem/item_all/internal/storage"
)

func TestWALAllowsReadWhileWriteTransactionIsOpen(t *testing.T) {
	store := openTestStore(t)
	ctx := context.Background()
	if _, err := store.ExecContext(ctx, `CREATE TABLE wal_values(value INTEGER NOT NULL) STRICT`); err != nil {
		t.Fatal(err)
	}
	if _, err := store.ExecContext(ctx, `INSERT INTO wal_values(value) VALUES(1)`); err != nil {
		t.Fatal(err)
	}
	writeStarted := make(chan struct{})
	releaseWrite := make(chan struct{})
	writeDone := make(chan error, 1)
	go func() {
		writeDone <- store.WithTx(ctx, func(txCtx context.Context, q storage.Querier) error {
			if _, err := q.ExecContext(txCtx, `INSERT INTO wal_values(value) VALUES(2)`); err != nil {
				return err
			}
			close(writeStarted)
			<-releaseWrite
			return nil
		})
	}()
	<-writeStarted
	readDone := make(chan struct{})
	var count int
	var readErr error
	go func() {
		readErr = store.QueryRowContext(ctx, `SELECT COUNT(*) FROM wal_values`).Scan(&count)
		close(readDone)
	}()
	select {
	case <-readDone:
	case <-time.After(time.Second):
		t.Fatal("WAL read blocked by open write transaction")
	}
	if readErr != nil {
		t.Fatal(readErr)
	}
	if count != 1 {
		t.Fatalf("uncommitted row visible, count = %d", count)
	}
	close(releaseWrite)
	if err := <-writeDone; err != nil {
		t.Fatal(err)
	}
}

func TestStoreSerializesApplicationWrites(t *testing.T) {
	store := openTestStore(t)
	ctx := context.Background()
	if _, err := store.ExecContext(ctx, `CREATE TABLE serial_values(value INTEGER NOT NULL) STRICT`); err != nil {
		t.Fatal(err)
	}
	firstStarted := make(chan struct{})
	releaseFirst := make(chan struct{})
	firstDone := make(chan error, 1)
	go func() {
		firstDone <- store.WithTx(ctx, func(txCtx context.Context, q storage.Querier) error {
			if _, err := q.ExecContext(txCtx, `INSERT INTO serial_values(value) VALUES(1)`); err != nil {
				return err
			}
			close(firstStarted)
			<-releaseFirst
			return nil
		})
	}()
	<-firstStarted
	secondDone := make(chan error, 1)
	go func() {
		_, err := store.ExecContext(ctx, `INSERT INTO serial_values(value) VALUES(2)`)
		secondDone <- err
	}()
	select {
	case err := <-secondDone:
		t.Fatalf("second write escaped serialization: %v", err)
	case <-time.After(100 * time.Millisecond):
	}
	close(releaseFirst)
	if err := <-firstDone; err != nil {
		t.Fatal(err)
	}
	if err := <-secondDone; err != nil {
		t.Fatal(err)
	}
}

func TestIndependentStoreMapsLockTimeoutToBusy(t *testing.T) {
	path := filepath.Join(t.TempDir(), "itemall.db")
	first, err := Open(context.Background(), DefaultConfig(path))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = first.Close() })
	secondConfig := DefaultConfig(path)
	secondConfig.BusyTimeout = 200 * time.Millisecond
	second, err := Open(context.Background(), secondConfig)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = second.Close() })
	ctx := context.Background()
	if _, err := first.ExecContext(ctx, `CREATE TABLE lock_values(value INTEGER NOT NULL) STRICT`); err != nil {
		t.Fatal(err)
	}
	lockHeld := make(chan struct{})
	release := make(chan struct{})
	firstDone := make(chan error, 1)
	go func() {
		firstDone <- first.WithTx(ctx, func(txCtx context.Context, q storage.Querier) error {
			if _, err := q.ExecContext(txCtx, `INSERT INTO lock_values(value) VALUES(1)`); err != nil {
				return err
			}
			close(lockHeld)
			<-release
			return nil
		})
	}()
	<-lockHeld
	started := time.Now()
	_, err = second.ExecContext(ctx, `INSERT INTO lock_values(value) VALUES(2)`)
	elapsed := time.Since(started)
	if !storage.IsCode(err, storage.CodeBusy) {
		t.Fatalf("error = %v", err)
	}
	if elapsed < 150*time.Millisecond || elapsed > time.Second {
		t.Fatalf("busy timeout elapsed = %s", elapsed)
	}
	close(release)
	if err := <-firstDone; err != nil {
		t.Fatal(err)
	}
}
