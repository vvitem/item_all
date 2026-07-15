package sqlite

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/vvitem/item_all/internal/storage"
)

func TestOpenCreatesDatabaseAndAppliesMigration(t *testing.T) {
	path := filepath.Join(t.TempDir(), "data", "itemall.db")
	store, err := Open(context.Background(), DefaultConfig(path))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = store.Close() })
	absolute, err := filepath.Abs(path)
	if err != nil {
		t.Fatal(err)
	}
	if store.Path() != absolute {
		t.Fatalf("Path = %q, want %q", store.Path(), absolute)
	}
	if _, err := os.Stat(absolute); err != nil {
		t.Fatal(err)
	}
	var count int
	if err := store.QueryRowContext(context.Background(), `SELECT COUNT(*) FROM schema_migrations`).Scan(&count); err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Fatalf("migration count = %d", count)
	}
	var table string
	if err := store.QueryRowContext(context.Background(), `SELECT name FROM sqlite_schema WHERE type='table' AND name='storage_meta'`).Scan(&table); err != nil {
		t.Fatal(err)
	}
	if table != "storage_meta" {
		t.Fatalf("table = %q", table)
	}
}

func TestOpenRepeatedIsIdempotent(t *testing.T) {
	path := filepath.Join(t.TempDir(), "itemall.db")
	first, err := Open(context.Background(), DefaultConfig(path))
	if err != nil {
		t.Fatal(err)
	}
	if err := first.Close(); err != nil {
		t.Fatal(err)
	}
	second, err := Open(context.Background(), DefaultConfig(path))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = second.Close() })
	var count int
	if err := second.QueryRowContext(context.Background(), `SELECT COUNT(*) FROM schema_migrations`).Scan(&count); err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Fatalf("migration count = %d", count)
	}
}

func TestOpenEnforcesForeignKeys(t *testing.T) {
	store := openTestStore(t)
	ctx := context.Background()
	if _, err := store.ExecContext(ctx, `CREATE TABLE parent(id INTEGER PRIMARY KEY) STRICT`); err != nil {
		t.Fatal(err)
	}
	if _, err := store.ExecContext(ctx, `CREATE TABLE child(parent_id INTEGER NOT NULL REFERENCES parent(id)) STRICT`); err != nil {
		t.Fatal(err)
	}
	if _, err := store.ExecContext(ctx, `INSERT INTO child(parent_id) VALUES(404)`); err == nil {
		t.Fatal("foreign key violation accepted")
	}
}

func TestOpenVerifiesApprovedPragmasAcrossFourConnections(t *testing.T) {
	store := openTestStore(t)
	ctx := context.Background()
	connections := make([]interface{ Close() error }, 0, 4)
	for range 4 {
		conn, err := store.db.Conn(ctx)
		if err != nil {
			t.Fatal(err)
		}
		connections = append(connections, conn)
		var journal string
		var foreignKeys, timeout, synchronous int
		if err := conn.QueryRowContext(ctx, `PRAGMA journal_mode`).Scan(&journal); err != nil {
			t.Fatal(err)
		}
		if err := conn.QueryRowContext(ctx, `PRAGMA foreign_keys`).Scan(&foreignKeys); err != nil {
			t.Fatal(err)
		}
		if err := conn.QueryRowContext(ctx, `PRAGMA busy_timeout`).Scan(&timeout); err != nil {
			t.Fatal(err)
		}
		if err := conn.QueryRowContext(ctx, `PRAGMA synchronous`).Scan(&synchronous); err != nil {
			t.Fatal(err)
		}
		if !strings.EqualFold(journal, "wal") || foreignKeys != 1 || timeout != 5000 || synchronous != 1 {
			t.Fatalf("pragmas = %q/%d/%d/%d", journal, foreignKeys, timeout, synchronous)
		}
	}
	for _, conn := range connections {
		if err := conn.Close(); err != nil {
			t.Fatal(err)
		}
	}
}

func TestOpenClassifiesCorruptDatabase(t *testing.T) {
	path := filepath.Join(t.TempDir(), "itemall.db")
	if err := os.WriteFile(path, []byte("not a sqlite database"), 0o600); err != nil {
		t.Fatal(err)
	}
	store, err := Open(context.Background(), DefaultConfig(path))
	if store != nil {
		_ = store.Close()
		t.Fatal("partial Store returned")
	}
	if !storage.IsCode(err, storage.CodeCorrupt) {
		t.Fatalf("error = %v", err)
	}
}

func TestStoreRejectsUseAfterClose(t *testing.T) {
	store := openTestStore(t)
	if err := store.Close(); err != nil {
		t.Fatal(err)
	}
	if err := store.Close(); err != nil {
		t.Fatalf("second close = %v", err)
	}
	ctx := context.Background()
	if _, err := store.ExecContext(ctx, `CREATE TABLE forbidden(id INTEGER)`); !storage.IsCode(err, storage.CodeClosed) {
		t.Fatalf("Exec error = %v", err)
	}
	if _, err := store.QueryContext(ctx, `SELECT 1`); !storage.IsCode(err, storage.CodeClosed) {
		t.Fatalf("Query error = %v", err)
	}
	if err := store.QueryRowContext(ctx, `SELECT 1`).Scan(new(int)); !storage.IsCode(err, storage.CodeClosed) {
		t.Fatalf("QueryRow error = %v", err)
	}
}

func openTestStore(t *testing.T) *Store {
	t.Helper()
	store, err := Open(context.Background(), DefaultConfig(filepath.Join(t.TempDir(), "itemall.db")))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = store.Close() })
	return store
}
