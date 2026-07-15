package sqlite

import (
	"os"
	"path/filepath"
	"testing"
)

func TestPrepareDatabasePathCreatesParentWithoutTouchingDatabase(t *testing.T) {
	path := filepath.Join(t.TempDir(), "nested", "data", "itemall.db")
	got, err := prepareDatabasePath(path)
	if err != nil {
		t.Fatal(err)
	}
	want, err := filepath.Abs(path)
	if err != nil {
		t.Fatal(err)
	}
	if got != want {
		t.Fatalf("path = %q, want %q", got, want)
	}
	info, err := os.Stat(filepath.Dir(got))
	if err != nil {
		t.Fatal(err)
	}
	if !info.IsDir() {
		t.Fatal("parent is not directory")
	}
	if _, err := os.Stat(got); !os.IsNotExist(err) {
		t.Fatalf("database was touched: %v", err)
	}
	matches, err := filepath.Glob(filepath.Join(filepath.Dir(got), ".itemall-write-probe-*"))
	if err != nil {
		t.Fatal(err)
	}
	if len(matches) != 0 {
		t.Fatalf("write probe leaked: %v", matches)
	}
}
