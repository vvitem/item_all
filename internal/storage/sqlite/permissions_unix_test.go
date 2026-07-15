//go:build !windows

package sqlite

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestOpenSecuresDatabaseFileMode(t *testing.T) {
	path := filepath.Join(t.TempDir(), "itemall.db")
	store, err := Open(context.Background(), DefaultConfig(path))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = store.Close() })
	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if got := info.Mode().Perm(); got != 0o600 {
		t.Fatalf("mode = %o", got)
	}
}
